# Micro-service in go
## Part-10
## Simple File Upload

Let's take a look at how to upload files in Go using the HTTP package. One of the benefits of using Go is that the process is straightforward due to the fact that any HTTP handler has access to the HTTP request, which includes a stream of data called the request body.

The request body is a stream that is read using a reader, which is a very versatile interface that combines the **io.Closer** and **io.Reader** interfaces. This means that we can read the data gradually as it is sent, rather than buffering all the data at once. This approach is very efficient and allows us to control the number of bytes that are read, ensuring that we don't accept too much data.

To upload a file in Go, we can create an HTTP handler that receives a POST request and reads the file from the request body using the io.Reader interface.

In this repo, services are decomposed into smaller services. and the service we are working now is **product_image**. inside this we have **main.go**

here we first define some env variable

```go
import (
    "github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "/tmp/images", "Base path to save images")
```

Now we add the file storage code, to write to the local disk and this is our initial **main.go**

```go
func main() {
    env.Parse()
    l := hclog. New(
            &hclog. LoggerOptions {
            Name: "product-images",
            Level: hclog. LevelFromString(*logLevel),
        },
    )
    // create a logger for the server from the default logger
    sl = l.Standard Logger(&hclog. Standard LoggerOptions{InferLevels: true})
    // create the storage class, use local storage
    // max filesize 5MB
    stor, err := files. NewLocal (*basePath, 1024*1000*5)
    if err != nil {
        l.Error("Unable to create storage", "error", err)
        os. Exit (l)
    }
    // create the handlers
    fh = handlers. NewFiles (stor, l)
    // create a new serve mux and register the handlers
    sm = mux.NewRouter()
    // filename regex: {filename: [a-zA-Z]+\\. [a-z]{3}}
    // problem with FileServer is that it is dumb
    // create a new server
    s := http.Server{
        Addr: *bindAddress, // configure the bind address
        Handler: sm, // set the default handler
        ErrorLog: sl, // the logger for the server
        ReadTimeout: 5 * time. Second, // max time to read request from the client
        WriteTimeout: 10 time. Second. // max time to write response to the client
        IdleTimeout: 120 * time.Second, // max time for connections using TCP Keep-Alive
    }

    // start the server
    go func(){
        l.Info("Starting server", "bind_address", *bindAddress)

        err := s.ListenAndServe()
        if err != nil {
            l.Error("Unable to start server", "error", err)
            os.Exit(1)
        }
    }()

    // trap sigterm or interrupt and gracefully shutdown the server
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, os.Kill)

    // Block untill a signal is received
    sig := <-c
    l.Info("Shutting down server with", "signal", sig)

    // gracefully shutdown the server, waiting max 30 seconds for current operations to complete
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    s.Shutdown(ctx)
}
```

Now we will add file upload functionality to it using mux subrouter.

To define a handler for uploading images in Go using the Gorilla framework, we will create a sub-router and use the HandleFunc function to specify the path of our endpoint.

First, we define our handler and register it with the HTTP method POST using the Gorilla router. We will name our endpoint images, as we want to upload images for a product ID. We will use a regular expression to validate the ID parameter as a number between 0 and 9, with one or more digits.Since we cannot get the file name from the HTTP request without using multi-part request, we will pass it as a parameter in the URL.
```go
ph := sm.Get(http.MethodPost).Subrouter()
ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)
```
in the above code we pass the filename with regex expression ```/{filename:[a-zA-Z]+\\.[a-z]{3}}```

Now in the handler, we don't need to validate the filename, if it is not null it will satisfy the filename parameter.

```go
package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-images/files"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters

	f.saveFile(id, fn, rw, r)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
```

Till now, we work on how to save the file, now we will look on how to get file back.
```go
// get files
gh := sm.Methods(http.MethodGet).Subrouter()
gh.Handle(
    "/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
    http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
)
```
here we define new handler and add the same path as we saved the image. and the function we are using is **http.FileServer** and we are using **http.StripPrefix** to remove *images* because it is not in the file path

