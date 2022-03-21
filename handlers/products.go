package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		p.addProducts(w, r)
		return
	}
	//catch all
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		//int id must exists in URI request, get it only with SL, using a regex
		reg := regexp.MustCompile(`([0-9]+)`) //return the regexp object
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("invalid URI, more than one id")

			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("invalid URI, more than one capture group")

			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		//the last section of URL in a post request mus containg the id
		idString := g[0][1]
		//must be cast as int to query datastore\
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("invalid URI, unable to convert tu number", idString)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		//p.l.Println("got id", id)
		p.updateProducts(id, w, r)
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

//addProducts, internal handler func to add a new product data to  datastore
func (p *Products) addProducts(w http.ResponseWriter, r *http.Request) {
	// Products logger
	p.l.Println("Handle POST Products")
	//product type variable
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}
func (p *Products) updateProducts(id int, w http.ResponseWriter, r *http.Request) {
	// Products logger
	p.l.Println("Handle PUT Products")
	//product type variable
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	}
	//call updateProduct function in data
	err = data.UpdateProduct(id, prod)
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
