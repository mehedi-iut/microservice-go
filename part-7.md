# Micro-service in go
## Part-6
## API Documentation using Swagger

In this documentation, we will explore API Documentation using swagger in go using **go-swagger**[link](https://goswagger.io/)

To install **go-swagger** we will use **Makefile**. to use **Makefile** we need to install **make** in our system. This **Makefile** will also generate swagger yaml file
```make
check_install:
    which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
    GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
```
Then we need to run ```make swagger``` from the terminal

Before runnting the above **Makefile** we need to add the swagger documentation in our code

In the **products.go** in handlers, we need to add some swagger documentation so that when we run the **Makefile** it can generate the swagger yaml file
```go
// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"microservice-go/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products{
	return &Products{l}
}


func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil{
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct {}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler{
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil{
			p.l.Println("[ERROR] vallidating proudct", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
```

*N.B* There shouldn't be any space between swagger documentation and **package handlers** line. otherwise spec will not generate. but when we documentating the function of API there should be space between swagger documentation and function code.

Now, if we run ```make swagger`` it will generate **swagger.yaml** in the root folder where **Makefile** is.

if the swagger installation failed then you can follow this [linke](https://github.com/go-swagger/go-swagger/blob/master/docs/install.md#debian-packages-)

### Code changed
From previous episode the code has changed and inside the products handler, there are different go files to handle the different REST API method
the below one is **get.go**
```go
package handlers

import (
	"net/http"
	"product-api/data"
)

// getProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
```

Now we will add basic swagger documentation in our get function

```go
package handlers

import (
	"net/http"
	"product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products

// getProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
```

Now, we want to add response to our swagger documentation of get method, so we need to create **struct** in **products.go** in handler section with **swagger:response** tag
```go
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"microservice-go/data"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
        err = prod.Validate()
        if err != nil{
            p.l.Println("[ERROR] validating product", err)
            http.Error(rw,
            fmt.Sprintf("Error validating product: %s", err),
            http.StatusBadRequest,
            )
            return
        }

		// add the product to the context
		ctx := context.WithValue(r.Context, KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
```

Now we will use docs handler to have nice ui of swagger

