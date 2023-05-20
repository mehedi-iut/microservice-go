package handlers

import (
	"fmt"
	"net/http"
	"practice/data"
)

// swagger:route POST /products products createProduct
// Create a new product
// Responses:
//      200: productResponse
//		422: errorValidation
//		501: errorResponse

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
