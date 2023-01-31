# Micro-service in go
## Part-3
## RESTful services

In this doc, we will develop RESTful services using go standard library.
here we will use simple **struct** to get data. later on we will use database in place of **struct**. Now, create **data** directory and inside of data will have **products.go** which will contain the data using **Product** struct

```go
package data

import "time"

// Product defines the structure for an API product
type Product struct {
    ID int
    Name string
    Description string
    Price float32
    SKU string
    CreatedOn string
    UpdatedOn string
    DeletedOn string
}

var productList = []*Product{
    &Product{
        ID: 1,
        Name: "Latte",
        Description: "Frothy milky coffee",
        Price: 2.45,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),

    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short and strong coffee without milk",
        Price: 1.99,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },  
}
```
in the above code, we have simple product list which we will consume with API Call
Now we will create **Handler** to get the data and convert to **JSON**

```go
package handlers
import (
    "log"
    "net/http"
)

type Products struct {
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request){

}
```

this is our handler. Inside the ServeHTTP we need write code to return **ProductList**. so we need to look at **endcoding/json**. this standard go library convert the productList **struct** to **JSON**. First we register our product handler in serveMux in **main.go**

```go
package main

import (
    "context"
    "product-api/handlers"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func main(){
    l := log.New(os.Stdout, "product-api ", log.LstdFlags)
    ph := handlers.NewProducts(l)

    sm := http.NewServeMux()
    sm.Handle("/", ph)

    s := http.Server{
        Addr: ":9090",
        Handler: sm,
        ErrorLog: l,
        ReadTimeout: 5*time.Second,
        WriteTimeout: 10*time.Second,
        IdleTimeout: 120*time.Second,
    }

    go func(){
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

    sig := <- c
    log.Println("Got signal:", sig)

    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    s.Shutdown(ctx)
}
```
Here we only create our handler and register it to our ServeMux.

Now we want to return our information. so we need to add method in data package which contain the **productList**. whether we are using database or fixed length product list or text file. It is very useful to abstract the detail logic of where the data came from, from the rest of the code. rather than exposing the database connection and implementing the sql write operation in handlers, we put that in method of our data object. below we add **GetProducts** which will return the product list.

```go
package data

import "time"

// Product defines the structure for an API product
type Product struct {
    ID int
    Name string
    Description string
    Price float32
    SKU string
    CreatedOn string
    UpdatedOn string
    DeletedOn string
}

func GetProducts() []*Product {
    return productList
}

var productList = []*Product{
    &Product{
        ID: 1,
        Name: "Latte",
        Description: "Frothy milky coffee",
        Price: 2.45,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),

    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short and strong coffee without milk",
        Price: 1.99,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },  
}
```

Now in the handler, we use the **GetProducts** method and get the list of product

```go
package handlers
import (
    "log"
    "net/http"
    "product-api/data"
)

type Products struct {
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    lp := data.GetProducts()
}
```

Now, we need to convert lp to **JSON** array. we do that using **json/encoding** library. we will use **json.Marshal**

```go
package handlers
import (
    "log"
    "net/http"
    "product-api/data"
)

type Products struct {
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    lp := data.GetProducts()
    d, err := json.Marshal(lp)
    if err != nil{
        http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
    }
    rw.write(d)
}
```
we can use **struct tag** to modify the json output

```go
package data

import "time"

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}


func GetProducts() []*Product {
    return productList
}

var productList = []*Product{
    &Product{
        ID: 1,
        Name: "Latte",
        Description: "Frothy milky coffee",
        Price: 2.45,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),

    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short and strong coffee without milk",
        Price: 1.99,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },  
}
```

Now we will get the output accroding to the struct tag and struct tag with **-** will be omitted.

we can also use [Encode](https://golang.org/pkg/encoding/json/#NewEncoder) method of **json** package. it can directly write to **io.Writer**. we also know that **http.ResponseWriter** is also implement **io.Writer** interface. So we can directly write our data to **http.ResponseWriter**. It is faster than **Marshal** and it doesn't use **buffer** memory i.e it doesn't allocate the memory for json. but If we have a huge JSON file, then we can consider the **Marshal** method. while this is not necessary for single thread, but when we consider the microservice, these microservices will have a number of different threads, it's doing a number of different concurrent operations, we want to be able to take advantange of all of the performance advantages. It is a quick solution but not an over optimization. it just using a different method.

In the data package, we will declare the new struct for ProductList so that we can add method to the struct which will be responsible for converting productList to JSON. in this way we will keep our handlers clean

in **products.go**
```go
package data

import "time"

// Product defines the structure for an API product
type Product struct {
    ID int
    Name string
    Description string
    Price float32
    SKU string
    CreatedOn string
    UpdatedOn string
    DeletedOn string
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(p)
}

func GetProducts() []*Product {
    return productList
}

var productList = []*Product{
    &Product{
        ID: 1,
        Name: "Latte",
        Description: "Frothy milky coffee",
        Price: 2.45,
        SKU: "abc323",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),

    },
    &Product{
        ID: 2,
        Name: "Espresso",
        Description: "Short and strong coffee without milk",
        Price: 1.99,
        SKU: "fjd34",
        CreatedOn: time.Now().UTC().String(),
        UpdatedOn: time.Now().UTC().String(),
    },  
}
```

now in the handlers we will just call the **ToJSON** method to convert to **JSON** format

```go
package handlers
import (
    "log"
    "net/http"
    "product-api/data"
)

type Products struct {
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    lp := data.GetProducts()
    err := lp.ToJSON(rw)
    if err != nil{
        http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
    }
}
```
Now we implement the REST api get method. go standard library doesn't handle the rest api implementation very well. later we will use the framework **Gorilla**. there is another framework named **Gin**. 
When we hit the application, we registered our product handler in "/" path, it will call the **ServeHTTP** internally. So we need to implement the REST logic in the **ServeHTTP** method in our handler. before doing that we need to convert our **GetProducts** to internal method by changing the name to **getProducts** and call **GetProducts** from the internal method.

Now the handler code will look like
```go
package handlers
import (
    "log"
    "net/http"
    "product-api/data"
)

type Products struct {
    l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request){
    if r.Method == http.MethodGet{
        p.getProducts(rw, r)
        return
    }

    // for now, other method is skipped
    rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request){
    lp := data.GetProducts()
    err := lp.ToJSON(rw)
    if err != nil{
        http.Error(rw, "unable to encode json", http.StatusInternalServerError)
    }
}
```