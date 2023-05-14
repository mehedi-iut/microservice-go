package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"practice/data"
)

type Products struct {
	l  hclog.Logger
	p  *data.ProductModel
	pl *data.ProductInfo
}

func NewProducts(l hclog.Logger, p *data.ProductModel, pl *data.ProductInfo) *Products {
	return &Products{l, p, pl}
}

//func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		p.getProducts(rw, r)
//		return
//	}
//	if r.Method == http.MethodPost {
//		p.addProducts(rw, r)
//	}
//}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	// Fetch products from the database or perform any necessary operations
	p.l.Info("Handle GET Products")
	products, err := p.p.Latest()
	if err != nil {
		http.Error(rw, "Unable to retrieve products", http.StatusInternalServerError)
		return
	}

	// Set the response content type
	rw.Header().Set("Content-Type", "application/json")

	// Write the products as JSON response
	// err = json.NewEncoder(rw).Encode(products)
	for _, prod := range products {
		err = prod.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to encode JSON response", http.StatusInternalServerError)
			return
		}
	}

}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("Handle POST Product")

	//prod := &data.ProductInfo{}
	prod := r.Context().Value(KeyProduct{}).(*data.ProductInfo)
	//err := prod.FromJSON(r.Body)
	//if err != nil {
	//	http.Error(rw, "Unable to encode the JSON", http.StatusInternalServerError)
	//}

	//fmt.Println(prod)
	_, err := p.p.Insert(prod)
	if err != nil {
		p.l.Error("Unable to insert the data", "error", err)
		http.Error(rw, "Unable to add data to the database", http.StatusInternalServerError)
		return
	}

	// Send a response indicating success
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "Product added successfully")

}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	//	return
	//}

	p.l.Info("Handle PUT Product", name)
	prod := r.Context().Value(KeyProduct{}).(*data.ProductInfo)

	err := p.p.UpdateProductByName(name, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

	// Send a response indicating success
	rw.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(rw, "Product Modified successfully")
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.ProductInfo{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Info("Unable to deserializing product", "error", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
