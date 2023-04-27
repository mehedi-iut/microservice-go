# Microservice in go
## refactor-3

this is the last episode where we will refactor our code.

in the **products-api** handler **get.go** ```prod, err := data.GetProductByID(id)``` **GetProductByID** is not a function anymore in the data package rather it is an **object** now.

Now we will refactor **products-api** handler function **products.go**

first, we have
```go
type Products struct {
    l *log.Logger
    v *data.Validation
    cc protos.CurrencyClient
}
```

now we will change **CurrencyClient** to **ProductsDB** as the functionality is now present in ProductsDB and we will change the builtin log package to hclog
```go
type Products struct {
    l hclog.Logger
    v *data.Validation
    productDB *data.ProductsDB
}
```

as we change the **Products** struct, we need to change the **NewProducts**
```go
func NewProducts(l *log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
	return &Products{l, v, cc}
}
```

we need to change the ```*log.logger``` and ```CurrencyClient```
```go
func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}
```

Now we will refactor the **get.go** handler
first in **ListAll**,
```go
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")

	prods, err := p.productDB.GetProducts("") // need to supply currency, we will fix that later on
    if err != nil{
        rw.WriteHeader(http.StatusInternalServerError)
        data.ToJSON(&GenericError{Message: err.Error()}, rw)
        return
    }

	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to  serializing product", "error", err)
	}
}
```

Now we will refactor **ListSingle** function, it is mostly replace log with hclog
```go
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	p.l.Debug("Get record", "id", id) // refactor to hclog

	prod, err := p.productDB.GetProductByID(id, "") // refactor

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to fetching product", "error", err) // refactor to hclog

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Unable to fetch product", "error", err) //refactor to hclog

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serializing product", err)
	}
}

```

same for **delete.go** in products-api handler, we only change log 
```go
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("Deleting record", "id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Unable to delete record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Unable to delete record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
```

Now in **post.go**
```go
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Debug("Inserting product: %#v\n", prod)
	p.productDB.AddProduct(prod)
}
```

in the **put.go**
```go
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Updating record", "id", prod.ID)

	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Debug("Product not found", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
```

Now our final one **main.go**
first we need to change ```l := log.New(os.Stdout, "product-api", log.LstdFlags)``` to ```l := hclog.Default()```

now we need to create **NewProductsDB** instance
```go
db := data.NewProductsDB(cc, l)
ph := handlers.NewProducts(l, v, db)
```

Now in the below *http.Server* **ErrorLog** take standard logger but we are using *hclog*, so we need to convert that to standard logger
```go
s := http.Server{
    Addr:         *bindAddress,      // configure the bind address
    Handler:      ch(sm),            // set the default handler
    ErrorLog:     l,                 // set the logger for the server
    ReadTimeout:  5 * time.Second,   // max time to read the request from the client
    WriteTimeout: 10 * time.Second,  // max time to write response to the client
    IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
}
```
now we will change **ErrorLog** to Standard logger with hclog
```go
s := http.Server{
    Addr:         *bindAddress,      // configure the bind address
    Handler:      ch(sm),            // set the default handler
    ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),                 // set the logger for the server
    ReadTimeout:  5 * time.Second,   // max time to read the request from the client
    WriteTimeout: 10 * time.Second,  // max time to write response to the client
    IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
}
```

Now we can run our code, and hit ```curl localhost:9090/products``` to get the product list. but I also want to specifiy the base currency like USD and url should looks like ```curl localhost:9090/products?currency=USD```

so in the **main.go**, we can add another method in the mux httpGet
```go
getR := sm.Methods(http.MethodGet).Subrouter()
getR.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
getR.HandleFunc("/products", ph.ListAll)

getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")
getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
```

Now we can fetch the currency from the url in **ListAll** and **ListSingle**
```cur := r.URL.Query().Get("currency")```

```go
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	p.l.Debug("Get record", "id", id) // refactor to hclog

    cur := r.URL.Query().Get("currency")

	prod, err := p.productDB.GetProductByID(id, cur) // refactor

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to fetching product", "error", err) // refactor to hclog

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Unable to fetch product", "error", err) //refactor to hclog

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serializing product", err)
	}
}
```


```go
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")

    cur := r.URL.Query().Get("currency")

	prods, err := p.productDB.GetProducts(cur) // need to supply currency, we will fix that later on
    if err != nil{
        rw.WriteHeader(http.StatusInternalServerError)
        data.ToJSON(&GenericError{Message: err.Error()}, rw)
        return
    }

	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to  serializing product", "error", err)
	}
}

```

lastly, we need to update our swagger documentation of List products as it takes url query
so, in the **docs.go**

```go
// swagger:parameters listProducts listSingleProduct
type productQueryParam struct {
	// Currency used when returning the price of the product,
	// when not specified currency is returned in GBP.
	// in: query
	// required: false
	Currency string
}
```

after that run ```make swagger```
