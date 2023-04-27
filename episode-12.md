# Microservice in Go
## Gzip

In web development, gzip refers to a method of compressing files to reduce their size and make them quicker to transfer over the internet. Gzip is a popular file format for compression and decompression that is widely supported by web servers, browsers, and other web tools.

When a web server receives a request for a file, it checks whether the client (usually a web browser) supports gzip compression. If the client does support it, the server will compress the file using gzip before sending it to the client. The client then decompresses the file and displays it to the user.

Gzip is especially useful for compressing large files, such as HTML, CSS, and JavaScript files, which are commonly used in web development. By compressing these files, they can be transferred more quickly over the internet, reducing the amount of time it takes for a web page to load.

Most modern web browsers and web servers support gzip compression, and it is often enabled by default. However, it is still important to ensure that gzip compression is properly configured on your web server to ensure that your web pages load quickly and efficiently.


In our case Gzip handler works as a middleware. because, before sending the data to client we need to zip the component
which we do in the middleware section

in **handlers** folder, we will create **zip_middleware.go**
```go
package handlers

import (
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			rw.Write([]byte("hello"))
		}
	})
}
```
Here, we define **GzipHandler** struct for our handler. after that we are not using **ServeHTTP** because to act as a 
middleware, we need **next** because 
*In Go, middleware functions are used to intercept and manipulate HTTP requests and responses between a client and 
a server. Middleware functions typically take a handler function as input, manipulate the request and response objects, 
and then call the next handler in the chain.*

we accomplished with **GzipMiddleware**, now in **http.HandleFunc**, it takes ResponseWriter and Request, then it check
header has **Accept-Encoding** with value **gzip**, because **Accept-Encoding** can contains other values [link](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding#syntax)
```Accept-Encoding: deflate, gzip;q=1.0, *;q=0.5``` that's why we use **strings.Contains** method instead of **==**.
then we simply write **Hello** to the response writer

Now, in the **main.go** we first need to declare middleware handler, then in the get router, we need to use
middleware function using gorilla **Use** function
```go
mw := handlers.GzipHandler{}

// we have gh as get router 
gh.Use(mw.GzipMiddleware)
```
Now we can run the curl with **--compressed** argument

```curl -v localhost:9091/images/1/test.png --compressed -o out.png```

Now, that is working, we need change our code to send Gzip response.
if it is Gzip response, we will create gziped response else we just simply call **next.ServeHTTP**

```go
package handlers

import (
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			return
		}
		// if it can't accept gzip, we handle normally
		next.ServeHTTP(rw, r)
	})
}
```
Now we have **http.ResponseWriter** which is a *interface** which has different method [link](https://pkg.go.dev/net/http#ResponseWriter)

so basically anything we do, we can as long as we're kind of implementing this response writer
interface, we can create our own response writers.

Now we can create wrapper on top of **http.ResponseWriter**. so we can create our own response writer i.e Gzip
ResponseWriter. So to create that we need to first create **struct**

```go
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}
```

our **WrappedResponseWriter**, we have two fields, one is **rw** which is **http.ResponseWriter** and another is
**gw** which is ```*gzip.Writer```.

so the gzip writer is in the compress gzip package of go and again it's one of the standard library features 
so if we look at that compress gzip [link](https://pkg.go.dev/compress/gzip) we can see that actually it has a reader and a writer 
so it's just implementing the standard go reader writer interfaces. now the kind of the nice thing about this is 
wherever you have a reader or writer you can use this gzipped gzip to read and write and it will automatically decode or encode
the contents of the stream or the reader or the writer into a zipped format 
so we're gonna put a gzip writer on there. so let's create that sort of idiomatic approach and let's create a new response writer 
```go
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedResponseWriter{rw: rw, gw: gw}
}
```


so *gw* it's *gzip.NewWriter()* so writer takes an i/o writer and what also implements an i/o writer? 
*http.ResponseWriter* which is really nice so now what we have above response writer 
and we also have the gzipped response writer.

Now we can implement *http.ResponseWriter* interface method in our new gzip response writer and now we will
implement *Header* method, as we see in doc [link](https://pkg.go.dev/net/http#ResponseWriter)

```go
func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}
```
here we are just returning the original http.ResponseWriter header

now we will implement *Write* method
```go
func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}
```
he signature for that is you're gonna write a byte of information sorry a slice of bytes you're going to 
return the length of the the data that you've written and an error if everything went wrong. 
when we do our *Write*, we gonna just do our *wr.rw.Write()* 
so we are using the gzipped writer because if you remember the gzipped writer wraps the response writer 
so when we call *Write* function to write any data that we are writing out is now going to be gzipped Writer

we need to add **WriteHeader**
```go
func (wr *WrappedResponseWriter) WriteHeader (statuscode int) {
	wr.rw.WriteHeader(statuscode)
}
```
so, same as other method, using http.ResponseWriter *WriteHeader* method, we wrapped with out custom responsewriter *wr*

```go
func (wr *WrappedResponseWriter) Flush(){
	wr.gw.Flush()
	wr.gw.Close()
}
```

the flush method is just going to do things like write anything so flush anything which hasn't actually been sent out 
on the underlying strings streams

Now, What is *Flush* function in *http.ResponseWriter*
In Go, the `http.ResponseWriter` interface is used to write the HTTP response from a server. The `Flush()` function 
is a method of this interface, and it is used to flush any buffered data to the client.
When you write data to the `http.ResponseWriter`, the data is buffered before it is sent to the client. 
This buffering is done for performance reasons - it allows the server to send the data in larger chunks, 
which can be more efficient than sending small pieces of data.

However, there are situations where you may want to force the buffered data to be sent immediately. 
This is where the `Flush()` function comes in. When you call `Flush()`, any buffered data is immediately sent to the client. 
This can be useful if you want to ensure that the client receives some data before the request is complete, 
or if you want to provide progress updates to the client during a long-running operation.
It's important to note that calling `Flush()` doesn't necessarily mean that the entire response has been sent to the client. 
There may still be more data to send, and the `Flush()` function just ensures that any buffered data is sent immediately.


Now, we have our own gzip ResponseWriter, now we will complete the `GzipMiddleware`
```go
func(g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip"){
			// create a gzipped response
			wrw := NewWrappedResponseWriter(rw)
			
			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
        }
		// handle if not gzip used
		next.ServeHTTP(rw, r)
    })
}
```