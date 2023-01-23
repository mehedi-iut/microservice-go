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
Here we import go package **net/http** which will be used handle http server related component

In the **main** function, we have used two method of http library **HandleFunc** and **ListenAndServe**

### HandleFunc
It register a function to a path on **defaultServeMux**. It tahke the function that is defined in its parameter and create a http *Handler* with it and added to the *defaultServeMux*

### defaultServeMux
It is a http **Handler**. Everything related to server in go is http **Handler**

### Handler
It is an *interface* in go http library. It has a method called ```
ServeHTTP(ResponseWriter, *Request)``` Any struct which has this method, implements the interface **Handler**

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