package data

import (
	"encoding/json"
	"io"
	"time"
)

//Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` //ignore it in response
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

//using a type for return the list of product allows to use methods on the structure (receivers)
//and encapsulates functionality
type Products []*Product

func (p *Products) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//ToJSON() use the json enconder method to return the list of product serialized as JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//REM: the http.responseWriter is a io.writer
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts returns a slice of products type
func GetProducts() Products {
	return productList
}
func AddProducts(p *Products) {

}
func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID
}

//productList hardcoded for now
var productList = []*Product{
	{
		ID:          1,
		Name:        "Late",
		Description: "Frothy mulky coffe",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffe without milk",
		Price:       1.99,
		SKU:         "fdj34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}