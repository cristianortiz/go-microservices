package data

import "testing"

//simple test to check the validate function
func TestCheckValidation(t *testing.T) {
	// product type,
	p := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "abs", //bad format to validateSku function
	}

	//should fails the test if p is empty struct or contains wrong values
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}

}
