package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/cristianortiz/go-microservices/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//getProducts, internal handler func to get products from datastore
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
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

//addProducts, internal handler func to add a new product data to  datastore
func (p *Products) AddProducts(w http.ResponseWriter, r *http.Request) {
	// Products logger
	p.l.Println("Handle POST Products")

	//extract the  product data from the context defined in the middleware executed by htttp request
	// before this handler, also need to cast the returned interface into product data type
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

//UpdateProduct() is handler func to udpate the data a specific product with their id
func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	//obtain the id parameter from http request through  g-mux
	vars := mux.Vars(r)
	//convert to int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "unable to cast id to int", http.StatusBadRequest)
		return
	}
	// Products logger
	p.l.Println("Handle PUT Products", id)

	//extract the  product data from the context defined in the middleware executed by htttp request
	// before this handler, also need to cast the returned interface into product data type
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	//call updateProduct function in data
	err = data.UpdateProduct(id, &prod)
	//ErrProductNotFound config in data.go
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

}

//type to store middleware response, in this case the product info
type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})

}
