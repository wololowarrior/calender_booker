package routes

import (
	"calendly_adventures/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")                                // Get the user
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")                                 // create user
	router.HandleFunc("/users/{id}/unavailable", handlers.CreateUnavailabilitySlots).Methods("POST") // create unavailability
	router.HandleFunc("/users/{id}/unavailable", handlers.GetUnavailabilitySlots).Methods("GET")     // get unavailability list

	router.HandleFunc("/users/{id}/event", handlers.CreateEvent).Methods("POST")                         // create event
	router.HandleFunc("/users/{id}/event/{event_id}", handlers.GetUnavailabilitySlots).Methods("DELETE") // create event
	router.HandleFunc("/users/{id}/event", handlers.GetEvents).Methods("GET")                            // list events
	router.HandleFunc("/users/{id}/meetings", handlers.GetMeetingsForAUser).Methods("GET")               // get reserved meetings

	router.HandleFunc("/meetings", handlers.CreateMeetings).Methods("POST")       // create meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.UpdateMeeting).Methods("PUT")    // reschedule meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.DeleteMeeting).Methods("DELETE") // cancel meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.GetMeeting).Methods("GET")       // get meeting

	router.HandleFunc("/meetings", handlers.GetMeetingsFromEvent).Methods("GET") // get available slots for a user

	return router
}
