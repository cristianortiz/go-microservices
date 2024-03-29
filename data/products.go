package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

//Product defines the structure for an API product,  fields are validated
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"` //ignore it in response
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

//using a type for return the list of product allows to use methods on the structure (receivers)
//and encapsulates functionality
type Products []*Product

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//Validate() use validator package to create a validator instance and applied to a product tpe
func (p *Product) Validate() error {
	//create a validator
	validate := validator.New()
	//register custom validation func
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)

}

//custom validation function for SKU product type field
func validateSKU(fl validator.FieldLevel) bool {
	//sku is of format abc-absd-sdfsd
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	//findallString returns a slice of byte
	matches := re.FindAllString(fl.Field().String(), -1)
	//must be only one pattern match
	if len(matches) != 1 {
		return false

	}
	return true

}

//ToJSON() use the json encoder method to return the list of product serialized as JSON
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

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)

}
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	//ensure the Id is correct
	p.ID = id
	//udpate the info of product p in position[pos] in datastore
	productList[pos] = p
	return nil

}

var ErrProductNotFound = fmt.Errorf("Product not found")

//findProduct() internal function to get the data of a product based on their id
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound

}

//getNextID assign an id when adding a new product to datastore
func getNextID() int {
	//get the last added product id
	lp := productList[len(productList)-1]
	return lp.ID + 1
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
