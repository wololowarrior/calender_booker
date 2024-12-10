package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	events, err := dao.GetAllEvents(userID)
	if err != nil {
		log.Printf("Error getting events: %s", err)
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
