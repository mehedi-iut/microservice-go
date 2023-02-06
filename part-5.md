# Micro-service in go
## Part-5

In this section we will use **gorilla** framework. previously we have below code
```go
package main

import (
	"context"
	"log"
	"micro-service-in-go/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main(){
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr: ":9090",
		Handler: sm,
		ErrorLog: l,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil{
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
```

First we will change our **http.NewServeMux** to gorilla **mux.NewRouter** which is our root router. then we can add **SubRouter** and in the **SubRouter** we can registered the handler function with path for specific http verb
```go
package main

import (
	"context"
	"log"
	"micro-service-in-go/handlers"
    "github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main(){
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()
    getRouter := sm.Methods("GET").Subrouter()
    getRouter.HandleFunc("/", ph.GetProducts)

	s := http.Server{
		Addr: ":9090",
		Handler: sm,
		ErrorLog: l,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil{
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
```

Now inside my handler function, we can delete the **ServeHTTP** function, we don't need that anymore, then we need to convert the private function to public ones by make it Capital.
```go
package handlers

import (
	"log"
	"micro-service-in-go/data"
	"net/http"
	"regexp"
	"strconv"
)


type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}


func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	err := lp.ToJSON(rw)
	if err != nil{
		http.Error(rw, "Unable to encode the JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request){
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil{
		http.Error(rw, "Unable to add the product", http.StatusBadRequest)
	}
	p.l.Printf("Prod %#v", prod)
	data.AddProduct(prod)
}

func(p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request){
	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to add the product", http.StatusBadRequest)
	}

	data.UpdateProduct(id, prod)
}
```

The above process summary is sets up routing for HTTP requests using the **mux** package in Go. The **mux.NewRouter()** function creates a new router and assigns it to the **sm** variable.

A sub-router is then created for HTTP **GET** requests using **sm.Methods("GET").Subrouter()** and assigned to the **getRouter** variable. The **HandleFunc** method is then used to register a function **ph.GetProducts** that will be executed when a GET request is made to the root URL **"/"**.

Now we will add other two method to it
```go
putRouter := sm.Methods(http.MethodPut).Subrouter()
putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)

postRouter := sm.Methods(http.MethodPost).Subrouter()
postRouter.HandleFunc("/", ph.AddProudct)
```

Now we need to use **middleware** to process the data that is send to server. Now what is **Middleware**?
    A middleware is a piece of code that is executed before or after an HTTP request is handled by a route in a web application. It provides a way to modify the request and response, or to perform some processing, such as authentication, authorization, logging, or error handling.
    In other words, middleware acts as an intermediary between the incoming request and the handler function that handles the request. It can perform operations on the request and/or response before or after the request is handled, and it can also choose to short-circuit the request handling and return a response immediately, or pass the request to the next middleware or the handler function.
    Middleware can be used to add common functionality to an application, such as adding headers to responses, handling CORS, validating requests, etc. In Go, middleware is often implemented as a function that takes in a http.Handler and returns a http.Handler.


