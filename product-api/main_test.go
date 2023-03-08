package main

import (
	"fmt"
	"microservice-go/sdk/client"
	"microservice-go/sdk/client/products"
	"testing"
)

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