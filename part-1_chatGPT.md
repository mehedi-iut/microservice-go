# Go HTTP Server

This documentation will guide you on how to create a simple HTTP server in Go.
## Prerequisites

* Go programming language installed
* Basic understanding of Go programming language

## Steps to Create a Simple HTTP Server in Go

1. Import the necessary packages, in this case **net/http** and **log**.
```go
package main
import (
    "net/http"
    "log"
)
```

2. In the **main** function, use the **http.HandleFunc** method to register a function as an HTTP handler for the root URL ("/"). This function will be called whenever a request is received with a URL that matches the pattern.
```go
func main(){
    http.HandleFunc("/", func(http.ResponseWriter, *http.Request){
        log.Println("Hello World")
    })
```

3. Use the http.ListenAndServe method to start the HTTP server and listen on a specific port, in this case 9090. This method uses the default http.ServeMux, which is a multiplexer that matches the URL of an incoming request and calls the corresponding handler to handle the request.
```go
    http.ListenAndServe(":9090", nil)
}
```

## Additional Information
### http.HandleFunc

**http.HandleFunc** is a convenience wrapper around **http.Handle** and is used to register a function as an HTTP handler. The function should have the signature **func(http.ResponseWriter, *http.Request)**. This means that any function that has this signature can be used as an HTTP handler.
### http.Handle

**http.Handle** is used to register an HTTP handler with the default **http.ServeMux**. It takes an **http.Handler** interface as its first argument and a string as its second argument, which is the pattern of the URL that the handler should match.
### http.ServeMux

**http.ServeMux** is a multiplexer or a request router that matches the URL of an incoming request and calls the corresponding handler to handle the request. In Go, the **http.ServeMux** is a struct that implements the **http.Handler** interface. It maintains a list of registered handlers, each associated with a pattern of URL path. The package **http** provides a default **ServeMux** named **DefaultServeMux** which is used by the **http.ListenAndServe** and **http.ListenAndServeTLS** functions if no other **http.ServeMux** is provided.
### http.ListenAndServe

**http.ListenAndServe** is a method that starts an HTTP server and listens on a specific address and port. It takes two arguments, the first being the address and port to listen on, and the second being the **http.Handler** to handle incoming requests. If the second argument is **nil**, it uses the default **http.ServeMux**.
### Handling Requests

Once the server is running, it will handle incoming requests by matching the URL of the request to the registered patterns and calling the corresponding handler. In the example provided, the function **func(http.ResponseWriter, *http.Request)** is registered as the handler for the root URL ("/"). This function simply logs "Hello World" to the console. You can customize this function to handle different types of requests and provide appropriate responses.
### Conclusion

This is a basic example of how to create a simple HTTP server in Go using the **net/http** package. You can build upon this example and add more functionality to your server, such as routing, middleware, and more. Remember to import necessary packages, register handlers, and start the server using **http.ListenAndServe** method and you're good to go.