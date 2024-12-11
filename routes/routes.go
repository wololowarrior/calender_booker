package routes

import (
	"calendly_adventures/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")                                // Get the user
	router.HandleFunc("/user", handlers.CreateUser).Methods("POST")                                 // create user
	router.HandleFunc("/user/{id}/unavailable", handlers.CreateUnavailabilitySlots).Methods("POST") // create unavailability
	router.HandleFunc("/user/{id}/unavailable", handlers.GetUnavailabilitySlots).Methods("GET")     // get unavailability list

	router.HandleFunc("/user/{id}/event", handlers.CreateEvent).Methods("POST")              // create event
	router.HandleFunc("/user/{id}/event/{event_id}", handlers.DeleteEvent).Methods("DELETE") // delete event
	router.HandleFunc("/user/{id}/event", handlers.GetEvents).Methods("GET")                 // list events
	router.HandleFunc("/user/{id}/meetings", handlers.GetMeetingsForAUser).Methods("GET")    // get reserved meetings
	router.HandleFunc("/user/{id}/overview", handlers.Overview).Methods("GET")               // get the days overview

	router.HandleFunc("/meetings", handlers.CreateMeetings).Methods("POST")       // create meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.UpdateMeeting).Methods("PUT")    // reschedule meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.DeleteMeeting).Methods("DELETE") // cancel meeting, done by others
	router.HandleFunc("/meetings/{id}", handlers.GetMeeting).Methods("GET")       // get meeting

	router.HandleFunc("/meetings", handlers.GetMeetingsFromEvent).Methods("GET") // get available slots for a user

	return router
}
