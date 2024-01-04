package routes

import (
	"github.com/gorilla/mux"
	"github.com/mukeshkuiry/eai-backend/handlers"
)

// Set up custom routes
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.LoadStockData).Methods("GET")
}
