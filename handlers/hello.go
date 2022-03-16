package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//hello is a simple handler
type Hello struct {
	l *log.Logger
}

//NewHello() creates a new hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle Hello requests")

	//read the request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Println("Error trying to get body request")
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	// write the response
	fmt.Fprintf(w, "Hello %s", b)
}
