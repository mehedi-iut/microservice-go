# Micro-service in go
## Part-1

Simple http server in go
```go
package main
import (
    "net/http"
    "log"
)

func main(){
    http.HandleFunc("/", func(http.ResponseWriter, *http.Request){
        log.Println("Hello World")
    })

    http.ListenAndServe(":9090", nil)
}
```
Here we import go package **net/http** which will be used to handle http server related component

In the **main** function, we have used two method of http library, **HandleFunc** and **ListenAndServe**

### HandleFunc
It register a function to a path on **defaultServeMux**. It takes the function that is defined in its parameter and create a http *Handler* with it and added to the *defaultServeMux*
So, **http.HandleFunc** is to register a function as an HTTP handler, it is a convenience wrapper around **http.Handle**. It is typically used to handle incoming HTTP requests and provide a response.
```go
http.HandleFunc("/", myHandler)
```

### http.Handle
In Go, **http.Handle** is a function that is used to register an HTTP handler with the default **http.ServeMux**. **http.ServeMux** is a multiplexer or a request router that matches the URL of an incoming request and calls the corresponding handler to handle the request.

**http.Handle** takes an **http.Handler** interface as its first argument and a string as its second argument, which is the pattern of the URL that the handler should match. When a request is received, the **http.ServeMux** will compare the request's URL to the registered patterns, and call the corresponding handler if a match is found.

For example, the following code will register the **myHandler** function as an HTTP handler for the root URL ("/")
```go
http.Handle("/", http.HandlerFunc(myHandler))
```

Alternatively you can use **http.HandleFunc** to register a function as an HTTP handler, it is a convenience wrapper around **http.Handle**
```go
http.HandleFunc("/", myHandler)
```
Once the handler is registered, the **http.ServeMux** will call it whenever a request is received with a URL that matches the pattern.

When you call **http.ListenAndServe** it will use the default ServeMux which is **http.DefaultServeMux** if you don't provide any **http.ServeMux** to it.


### defaultServeMux
It is a http **Handler**. Everything related to server in go is http **Handler**.
The **http.ServeMux** is a multiplexer or a request router that matches the URL of an incoming request and calls the corresponding handler to handle the request. In Go, the **http.ServeMux** is a struct that implements the **http.Handler** *interface*. It maintains a list of registered handlers, each associated with a pattern of URL path.

In Go, the package **http** provides a default **ServeMux** named **DefaultServeMux** which is used by the **http.ListenAndServe** and **http.ListenAndServeTLS** functions if no other **http.ServeMux** is provided.

You can register handlers with the default ServeMux by using **http.Handle** and **http.HandleFunc** functions.

### Handler
It is an *interface* in go http library. It has a method called ```
ServeHTTP(ResponseWriter, *Request)``` Any struct which has this method, implements the interface **Handler**. So HTTP handler is a function that implements the **http.Handler** interface. This interface is defined as:
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```
This means that any function that has the signature **func(http.ResponseWriter, *http.Request)** can be used as an HTTP handler. These functions are typically used to handle incoming HTTP requests and provide a response. For example, a simple HTTP handler might look like this:
```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}
```
This handler would write the string "Hello, World!" to the response when it is called.

### ListenAndServe
It is a method of **net/http**. It construct http server and register a default **Handler** to it. If we don't specify any **Handler** to it, it will use **defaultServeMux**

### ServeMux
It is a http *Handler*. It is responsible for redirecting path i.e, you map a function at a path and when request come to server, **ServeMux** will determine which function need to execute based on path. if we don't setup anything using *ServeMux*, it will use *defaultServeMux*. It contains a function called **Handle**. It is responsible for register **Handler** to a path

```go
package main
import (
    "net/http"
    "log"
)

func main(){
    http.HandleFunc("/", func(http.ResponseWriter, *http.Request){
        log.Println("Hello World")
    })

    http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request){
        log.Println("Goodbye World")
    })

    http.ListenAndServe(":9090", nil)
}
```
In the above code there is two path. if we hit other than two path let's say "/api", it will fall back to "/" and will show "Hello World" message

### Response Writer
It is an interface used by http **Handler** to construct http response. It has number of method. It can write to client. So we can use it to show message to client. If we want to send back status error, we can use **WriteHeader** method

```go
package main
import (
    "net/http"
    "log"
    "io/ioutil"
)

func main(){
    http.HandleFunc("/", func(rw http.ResponseWriter, r*http.Request){
        log.Println("Hello World")
        d, err := ioutil.ReadAll(r.Body)
        if err != nil{
            rw.WriteHeader(http.StatusBadRequest)
            rw.Write([]byte("Ooops"))
            return
        }

        fmt.Fprintf(rw, "Hello %s", d)
    })

    http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request){
        log.Println("Goodbye World")
    })

    http.ListenAndServe(":9090", nil)
}
```

In the above code we write back to the client. If we get any error, we also write back to the client with error code. We can reduce the code inside the error block using **http.Error**

```go
package main
import (
    "net/http"
    "log"
    "io/ioutil"
)

func main(){
    http.HandleFunc("/", func(rw http.ResponseWriter, r*http.Request){
        log.Println("Hello World")
        d, err := ioutil.ReadAll(r.Body)
        if err != nil{
            http.Error(rw, "Ooops", http.StatusBadRequest)
            return
        }

        fmt.Fprintf(rw, "Hello %s", d)
    })

    http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request){
        log.Println("Goodbye World")
    })

    http.ListenAndServe(":9090", nil)
}
```