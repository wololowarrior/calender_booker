package routes

import (
	"calendly_adventures/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}/unavailable", handlers.CreateUnavailabilitySlots).Methods("POST")
	router.HandleFunc("/users/{id}/unavailable", handlers.GetUnavailabilitySlots).Methods("GET")
	return router
}
