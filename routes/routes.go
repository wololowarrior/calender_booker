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

	router.HandleFunc("/users/{id}/event", handlers.CreateEvent).Methods("POST")                         // create event
	router.HandleFunc("/users/{id}/event/{event_id}", handlers.GetUnavailabilitySlots).Methods("DELETE") // create event
	router.HandleFunc("/users/{id}/event", handlers.GetEvents).Methods("GET")                            // list events
	router.HandleFunc("/users/{id}/meetings?date", handlers.GetUnavailabilitySlots).Methods("GET")       // get reserved meetings
	router.HandleFunc("/meetings", handlers.GetUnavailabilitySlots).Methods("POST")                      // create meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.GetUnavailabilitySlots).Methods("PUT")                  // reschedule meeting, done by others
	router.HandleFunc("/meetings{id}", handlers.GetUnavailabilitySlots).Methods("DELETE")                // cancel meeting, done by others

	return router
}
