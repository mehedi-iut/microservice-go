# Microservice in go
## Part-2

Up untill now, what we have done is everyting inside **main.go**.
```go
package main
import (
    "net/http"
    "fmt"
    "log"
    "io/ioutil"
)

func main(){
    http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request){
        log.Println("Hello World")
        d, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(rw, "Oops", http.StatusBadRequest)
            return
        }

        fmt.Fprintf(rw, "Hello %s", d)
    })

    http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request){
        log.Println("Goodbye World")
    })

    http.ListenAndServer(":9090", nil)
}
```

We need to create **handlers** in separate module and import it in the **main.go**. To create a **Handler** we need to implement the *interface* **Handle** which has one method called **ServeHTTP**. so if we want to implement the *interface* **Handle** we need to add the method **ServeHTTP**. Below is the skeleton of our **Hello** handler

```go
package handlers
import (
    "net/http"
)

type Hello struct {

}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request){

}
```

Now we have **Hello** handler which can be used to register in the **ServeMux**

Now, we can copy the whole code inside **http.HandleFunc** function block

```go
package handlers
import (
    "net/http"
)

type Hello struct {

}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    log.Println("Hello World")
    d, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(rw, "Oops", http.StatusBadRequest)
        return
    }

    fmt.Fprintf(rw, "Hello %s", d)
}
```

Here, we need to change few things to follow the idomatic principal of go.
inside the handler, we should not use log, because we need more control over the log. And we shouldn't use any object inside the handler.

Now, inside our **Hello** handler, we need to define the **log.Logger**.
```go
type Hello struct {
    l *log.Logger
}
```
by following the idiomatic principal of go we will create a function which take logger and return the handler
```go
func NewHello(l *log.Logger) *Hello {
    return &Hello{l}
}
```

The purpose of the **NewHello** function is to create and initialize a new **Hello** struct. It is commonly known as a **factory function**, which is a **pattern** in Go that provides a convenient way to create objects with a specific set of dependencies.

The **NewHello** function takes one argument, a pointer to a **log.Logger**, which is a dependency that the **Hello** struct needs in order to function. By passing in the **log.Logger** as an argument, the **NewHello** function allows the caller to provide a specific implementation of the **log.Logger** interface, such as a custom logger or a third-party logger library.

By using the **factory pattern**, the **NewHello** function decouples the creation of the **Hello** struct from its implementation, making it more flexible and easier to test. For example, in a test case, you can pass a mock logger to the **NewHello** function and check the behavior of the Hello struct without actually logging any messages.

The NewHello function also makes the code more readable, as it clearly communicates the dependencies required to create a new **Hello** struct. This makes it easy to understand what the struct needs to function, and also makes it easy to change the dependencies if needed in future.

So the final **Hello** handler code
```go
package handlers
import (
    "net/http"
)

type Hello struct {
    l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
    return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    h.l.Println("Hello World")
    d, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(rw, "Oops", http.StatusBadRequest)
        return
    }

    fmt.Fprintf(rw, "Hello %s", d)
}
```

Now we need to change **main.go** to use our **Hello** handler
```go
package main

import (
    "fmt"
    "working/handlers"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)
func main(){
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hh := handlers.NewHello(l)

    http.ListenAndServe(":9090", nil)
}
```

here, we have **working** folder and inside that we have **handlers** folder where all the handler code reside

first we create New log instance. log can write to files as its first parameter is **io.Writer**. but here we use **os.Stdout** as io.Writer, then we have prefix and flag parameter. then we create handlers instance **hh** with **NewHello** with log **l** as Hello handler depends on log


*Note*: When a request is sent to a server, the server has a default handler, which is typically an instance of **http.ServeMux**. **http.ServeMux** is a multiplexer, or a "request router", that is responsible for determining which handler should handle the request based on the request's URL path.

When the server receives a request, it calls the **ServeHTTP** method on the default handler (**http.ServeMux**) with the request and response writer as arguments. The **ServeHTTP** method has special logic that looks at the request's URL path and compares it to a list of registered handlers and their corresponding URL patterns. If it finds a match, it calls the ServeHTTP method of the corresponding handler, passing the request and response writer as arguments.

The registered handlers are typically instances of the **http.Handler** interface, which has a single method called **ServeHTTP**. This method is responsible for handling the request and writing the response.

Here is a summary of the process:

1. A client sends a request to the server.
2. The server receives the request and calls the ServeHTTP method on its default handler, which is an instance of http.ServeMux.
3. http.ServeMux compares the request's URL path to a list of registered handlers and their corresponding URL patterns.
4. If it finds a match, it calls the ServeHTTP method of the corresponding handler, passing the request and response writer as arguments.
5. The registered handler's ServeHTTP method handles the request, and writes the response to the response writer.

It's important to notice that the ServeMux does not handle the request itself, it just route it to the appropriate handler based on the request's URL path. The handler is responsible for handling the request and writing the response.

The code uses the **http.ListenAndServe** function to start an HTTP server. This function takes two arguments: an address (in the form of an IP address and a port number), and an HTTP handler. In this case, the address provided is **:9090** which is a shorthand for *any IP address on port 9090*. The second argument is a variable **sm** which is a ServeMux, it is a multiplexer for HTTP requests.

By default, if no handler is provided as the second argument, the **http.ListenAndServe** function will use the default ServeMux. However, in this case, we want to use our own ServeMux instead of the default one. We create a new ServeMux by calling the **http.NewServeMux()** function, and register a handler for the root path ("/") to it by calling the **sm.Handle("/", hh)** method.

Finally, we pass our custom ServeMux to the **http.ListenAndServe** function, so that it uses our ServeMux to handle all incoming requests.

In short, the code creates an HTTP server on port 9090 and uses a custom ServeMux, which is configured to handle requests for the root path ("/") with a specific handler. This handler is created using the "handlers" package and passed a log object for logging.
```go
package main

import (
    "fmt"
    "working/handlers"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)
func main(){
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hh := handlers.NewHello(l)

    sm := http.NewServeMux()
    sm.Handle("/", hh)

    http.ListenAndServe(":9090", sm)
}
```

Now we will add our second handler i.e goodbye handler. create **goodbye.go** under **handlers** and add the below code.
```go
package handlers
import (
    "log"
    "net/http"
)

type Goodbye struct {
    l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
    return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    rw.Write([]byte("Byee"))
}
```

we need to add that in our **main.go**
```go
package main

import (
    "fmt"
    "working/handlers"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)
func main(){
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hh := handlers.NewHello(l)
    gh := handlers.NewGoodbye(l)

    sm := http.NewServeMux()
    sm.Handle("/", hh)
    sm.Handle("/goodbye", gh)

    http.ListenAndServe(":9090", sm)
}
```

The timeout in Golang is important because it helps manage the finite resources of a server. When a client connects to a server and performs an action, it creates a blocked connection. If there are too many blocked connections, the server will stop responding to new requests. To prevent this, we can fine-tune the server by creating a new [http.Server](https://pkg.go.dev/net/http#Server). In the documentation, there are options for adjusting the ReadTimeout, WriteTimeout, and IdleTimeout. The ReadTimeout can be set to a large or small value depending on the size of the file being read. The WriteTimeout depends on the amount of data being sent. The IdleTimeout manages the connection duration. Keeping the connection open with the same client can improve performance by reducing the time needed for DNS queries, TCP handshakes, and other processes. This is especially useful when there are many micro-services connecting to each other. When using TLS, which is a bit more time-consuming to establish a connection, it is important to keep the connection open. For random clients, the IdleTimeout should be set to a low value.

```go
package main

import (
    "fmt"
    "working/handlers"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)
func main(){
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hh := handlers.NewHello(l)
    gh := handlers.NewGoodbye(l)

    sm := http.NewServeMux()
    sm.Handle("/", hh)
    sm.Handle("/goodbye", gh)

    s := &http.Server{
        Addr: ":9090",
        Handler: sm,
        IdleTimeout: 120*time.Second,
        ReadTimeout: 1*time.Second,
        WriteTimeout: 1*time.Second,
    }

    s.ListenAndServer()
}
```


Now we need to add the **Graceful Shutdown**. Graceful shutdown is a process in which a server stops accepting new connections and waits for existing connections to complete before shutting down. This ensures that ongoing requests are not disrupted and allows for a clean and orderly termination of the server, avoiding any loss of data or incomplete requests. A graceful shutdown allows the server to complete any remaining work, free up resources, and ensure that the shutdown is performed in a controlled manner.

In go server we can use **Shutdown**[link](https://pkg.go.dev/net/http#Server.Shutdown) function. it accepts two parameters, **context** and **timeout**. It is used for graceful shutdown.


```go
package main

import (
    "fmt"
    "working/handlers"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)
func main(){
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hh := handlers.NewHello(l)
    gh := handlers.NewGoodbye(l)

    sm := http.NewServeMux()
    sm.Handle("/", hh)
    sm.Handle("/goodbye", gh)

    s := &http.Server{
        Addr: ":9090",
        Handler: sm,
        IdleTimeout: 120*time.Second,
        ReadTimeout: 1*time.Second,
        WriteTimeout: 1*time.Second,
    }

    go func(){
        err := s.ListenAndServe()
        if err != nil {
            l.Fatal(err)
        }
    }()

    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, os.Interrupt)
    signal.Notify(sigChan, os.Kill)

    sig := <-sigChan
    l.Println("Recived terminate, graceful shutdown", sig)

    tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
    s.Shutdown(tc)
}
```

This code is for starting an HTTP server in Go and handling a graceful shutdown.

* The first section starts a goroutine (concurrent function) which calls the ListenAndServe method on the "s" server object. If there is an error, it is logged using the "l" logger.

* The second section sets up a channel, "sigChan", to receive system signals. The function "signal.Notify" is called twice to register this channel to receive either an "Interrupt" or "Kill" signal.

* The third section waits for a signal to be received on the "sigChan" channel, and logs the received signal.

* The final section creates a timeout context, "tc", with a timeout of 30 seconds using the "WithTimeout" function from the "context" package. The server's "Shutdown" method is then called with this context, allowing it to initiate a graceful shutdown by closing any new incoming connections and allowing existing connections to complete.

*Why goroutine is needed in this case?*

A goroutine is used in this case to run the server in a concurrent manner, allowing it to handle multiple requests in parallel without blocking the main execution flow of the program.

Using a goroutine in this case ensures that the server can continue running and handling requests while the main function continues to execute. This is important because the main function must wait for a signal to initiate a graceful shutdown, and it would not be able to do so if it were blocked by the server.

By running the server in a separate goroutine, the main function can continue to run and handle the shutdown process while the server continues to operate in the background. This allows for a clean and efficient handling of the server shutdown, without interrupting ongoing requests or causing any data loss.

```go
sig := <-sigChan
```
using this code block, we pause the execution of main program and the below code will not execute untill the we recieve any signal in our channel. when we recieve the signal from the channel, then the rest of the code will execute