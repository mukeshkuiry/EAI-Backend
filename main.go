package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mukeshkuiry/eai-backend/handlers"
	"github.com/mukeshkuiry/eai-backend/routes"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Set up custom routes
	routes.SetupRoutes(r)

	handlers.FetchAndStoreStock()

	// Start the server
	port := 4000
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
