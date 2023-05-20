package handlers

import (
	"github.com/hashicorp/go-hclog"
	"practice/data"
)

type Products struct {
	l  hclog.Logger
	p  *data.ProductModel
	pl *data.ProductInfo
	v  *data.Validation
}

func NewProducts(l hclog.Logger, p *data.ProductModel, pl *data.ProductInfo, v *data.Validation) *Products {
	return &Products{l, p, pl, v}
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

type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

type KeyProduct struct{}
