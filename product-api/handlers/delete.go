package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"practice/data"
)

// swagger:route DELETE /products/{name} products deleteProduct
// Update a list of products
// Responses:
//		201: noContentResponse
//  	404: errorResponse
//  	501: errorResponse

func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	p.l.Info("Deleting Product", name)
	rw.Header().Add("Content-Type", "application/json")

	err := p.p.DeleteProductByName(name)
	if err == data.ErrProductNotFound {
		p.l.Error("Deleting record name doesn't exist", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Deleting record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
