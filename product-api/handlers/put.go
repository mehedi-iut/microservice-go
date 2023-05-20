package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"practice/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// Responses:
//		201: noContentResponse
//		404: errorResponse
//		422: errorValidation

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//name := vars["name"]
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	//	return
	//}

	name := chi.URLParam(r, "name")

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
