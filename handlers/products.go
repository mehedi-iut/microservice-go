package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practice/data"
	"github.com/hashicorp/go-hclog"
)

type Products struct {
	l  hclog.Logger
	p  *data.ProductModel
	pl *data.ProductInfo
}

func NewProducts(l hclog.Logger, p *data.ProductModel, pl *data.ProductInfo) *Products {
	return &Products{l, p, pl}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {

		p.addProducts(rw, r)
	}
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	// Fetch products from the database or perform any necessary operations
	products, err := p.p.Latest()
	if err != nil {
		http.Error(rw, "Unable to retrieve products", http.StatusInternalServerError)
		return
	}

	// Set the response content type
	rw.Header().Set("Content-Type", "application/json")

	// Write the products as JSON response
	err = json.NewEncoder(rw).Encode(products)
	if err != nil {
		http.Error(rw, "Unable to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	prod := &data.ProductInfo{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to encode the JSON", http.StatusInternalServerError)
	}

	fmt.Println(prod)
	_, err = p.p.Insert(prod)
	if err != nil {
		p.l.Error("Unable to insert the data", "error", err)
		http.Error(rw, "Unable to add data to the database", http.StatusInternalServerError)
		return
	}

}
