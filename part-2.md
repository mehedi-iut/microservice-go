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
