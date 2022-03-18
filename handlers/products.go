package handlers

import (
	"log"
	"net/http"

	"github.com/cristianortiz/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and satisfice the http.Handler // interface
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check the HTTP request method to call the apropiate handler func
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}
	//catch all
	//if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)

}

//getProducts, internal handler func to get products from datastore
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	// Products logger
	p.l.Println("Handle GET Products")
	// fetch the products from the datastore
	lp := data.GetProducts()
	//now using ToJSON the encoding and serialize the listo to JSON,
	// data and the response are encapsulated in one function
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "unable to marshal json", http.StatusInternalServerError)
	}

}

//addProducts, internal handler func to get products from datastore
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	// Products logger
	p.l.Println("Handle POST Products")
	product := &data.Products{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	}
}
