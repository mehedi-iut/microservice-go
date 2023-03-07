# Micro-service in go
## Part-8
## Swagger Client

From previous episode code, we can now add swagger client in our project. but why we need swagger client?

*the function of the Swagger client in Go is to automate the generation of client libraries for RESTful APIs based on the Swagger specification, making it easier for developers to interact with the APIs and reducing the amount of boilerplate code they need to write.*

Now to create swagger client first we need to create directory in root called **sdk**. 
```bash
mkdir sdk
cd sdk
```

Now we need to run swagger command to generate client
``` bash
swagger generate client -f ../microservice-go/swagger.yaml -A product-api
```

you can also use ```$(go env GOPATH)/bin/swagger``` to generate the client.

In the root along with **main.go** create a test file called **main_test.go** and add following code
```go
package main

import (
	"microservice-go/sdk/client"
	"microservice-go/sdk/client/products"
	"testing"
)

func TestOurClietn(t *testing.T){
	c := client.Default
	params := products.NewListProductsParams()
	c.Products.ListProducts(params)
}
```

if we go to the definition of **ListProducts** we can see it takes **params**. So we define **params** using **NewListProductsParams()**.

Now we want to print the response from **ListProducts**

```go
func TestOurClietn(t *testing.T){
	c := client.Default
	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println(prod)
}
```
if we run our code, we will get **connection refused** error, because it is listening in port 80 but our code is running in port 9090. It gets this default settings from swagger documentation when we declared it in **docs.go** file  with **swagger:meta**

but we can override the default settings

if we look into **sdk/client/product_api_client.go** we can see there is **NewHTTPClientWithConfig** and to use that we need to define config and pass it to the **NewHTTPClientWithConfig**

```go
func TestOurClietn(t *testing.T){
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
    c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println(prod)
}
```
if we run the test we will see, new client is working, but we get different error

*&[] (*[]*models.Product) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface*

Now we need to fix this issue. if we debug the test (you can check the video of this episode). we will find that our service return plain text but it should return application/json as we defined it in the swagger documentation

In order to fix that, we need to add header in our get method in **get.go**
```go
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

    rw.Header().Add("Content-Type", "application/json")

	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}
```

Now if we run the test, it will pass.
now if we want to see our product list we can add ```fmt.Printf("%#v", prod.GetPayload()[0])```
and we need to fail the test manually using ```t.Fail()```

```go
func TestOurClietn(t *testing.T){
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil{
		t.Fatal(err)
	}

	fmt.Printf("%#v",prod.GetPayload()[0])
	t.Fail()
}
```

Lastly, we need to add ```rw.Header().Add("Content-Type", "application/json")``` to **put.go** and **delete.go**


