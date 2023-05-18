package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"practice/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// Responses:
//     200: productResponse

func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
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
		err = data.ToJSON(prod, rw)
		if err != nil {
			http.Error(rw, "Unable to encode JSON response", http.StatusInternalServerError)
			return
		}
	}

}

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	p.l.Info("Get record name", name)

	product, err := p.p.GetProductByName(name)
	//p.l.Error("Can't get item", "Error", err)
	if err != nil {
		http.Error(rw, "Unable to get the item", http.StatusInternalServerError)
		return
	}

	err = data.ToJSON(product, rw)
	if err != nil {
		http.Error(rw, "Unable to encode JSON response", http.StatusInternalServerError)
		return
	}

}
