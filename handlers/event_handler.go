package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"calendly_adventures/dao"
	"calendly_adventures/models"
	"github.com/gorilla/mux"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'userid' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "'userid' must be an integer", http.StatusBadRequest)
		return
	}

	_, err = dao.GetUser(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Error retrieving user: %v", err)
		}
		return
	}

	decoder := json.NewDecoder(r.Body)
	var event models.Event
	err = decoder.Decode(&event)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	event.UID = userID
	err = dao.InsertEvent(&event)
	if err != nil {
		log.Printf("Error creating event: %s", err)
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'userid' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "'userid' must be an integer", http.StatusBadRequest)
		return
	}

	_, err = dao.GetUser(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Error retrieving user: %v", err)
		}
		return
	}

	events, err := dao.GetAllEvents(userID)
	if err != nil {
		log.Printf("Error getting events: %s", err)
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'userid' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "'userid' must be an integer", http.StatusBadRequest)
		return
	}

	_, err = dao.GetUser(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Error retrieving user: %v", err)
		}
		return
	}
	eventIDStr, ok := vars["event_id"]
	if !ok {
		http.Error(w, "missing 'eventId' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "'eventId' must be an integer", http.StatusBadRequest)
		return
	}
	_, err = dao.GetEvent(eventID)
	if err != nil {
		log.Print(err.Error())
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
	}
	err = dao.DeleteEvent(eventID, userID)
	if err != nil {
		log.Printf("Error deleting event: %s", err)
		if strings.Contains(err.Error(), "event doesn't exist") {
			http.Error(w, "event doesn't exist", http.StatusNotFound)
		}
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
