package handlers

import (
	"fmt"
	"net/http"
)

// HomeHandler is the handler for the home route "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, welcome to the home page!")
}
