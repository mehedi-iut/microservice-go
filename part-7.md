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

Now we will use docs handler to have nice ui of swagger. To create handler for it we will use **ReDoc**[link](https://github.com/Redocly/redoc)

To import **ReDoc**, 
```go
import (
	"github.com/go-openapi/runtime/middleware"
)
```

And to create the handler in the path **/docs**, we will use already define **getRouter**
```go
opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
sh := middleware.Redoc(opts, nil)

getRouter.Handle("/docs", sh)
```

Now if we run our server and visit **localhost:9090/docs** we will get error, because **ReDoc** try to download the **swagger.yaml** from our server but it can't find it, this is because, we don't serve the **swagger.yaml** file from our server

To serve **swagger.yaml** file, we need to add below code
```go
getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
```

we will add **delete** method in our API. so to do that, I will add **mux** router method in our **main.go**

```go
deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
deleteRouter.HandleFunc("/{id:[9-0]+}", ph.DeleteProduct)
```

we will create new **delete.go** in our **handlers**
```go
package handlers
import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"microservice-go/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
//  201: noContent

// DeleteProduct deletes a product from the database
func (p *Products) DeleteProduct (rw http.ResponseWriter, r *http.Request) {
	// this will always convert because of the router
	vars := mux.Vars(r)
	id, _ := strconv. Atoi(vars ["id"])

	p.l.Println("Handle DELETE Product", id)
	err := data.DeleteProduct(id)

	if err == data. ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
```

here we add also swagger documentation for our **DeleteProduct**. so it takes **id** parameter and in responses it has **noContent**. so we need to define those two parameters in **products.go**. so our swagger documentation is 
```
// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
//  201: noContent
```

In **products.go**, we need to define two struct for **id** parameter and **noContent** as we have used it in our swagger documentation.

```go
// swagger:response noContent
type productsNoContent struct {
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}
```

From the above code, we will find basic description in the swagger ui. but we can add rich description using swagger **model** tag.
So we can add swagger:model in the **data** handler product struct.
```go
// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this user
	//
	// required: true
	// min:1
	ID 			int 	`json:"id"`
	Name 		string 	`json:"name" validate:"required"`
	Description string 	`json:"description"`
	Price 		float32	`json:"price" validate:"gt=0"`
	SKU 		string 	`json:"sku" validate:"required,sku"`
	CreatedOn   string 	`json:"-"`
	UpdatedOn   string 	`json:"-"`
	DeletedOn   string 	`json:"-"`
}
```
