package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"calendly_adventures/dao"
	"calendly_adventures/models"
	"github.com/gorilla/mux"
)

func CreateUnavailabilitySlots(w http.ResponseWriter, r *http.Request) {
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
	var input models.UnavailableSlots
	err = decoder.Decode(&input)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	input.UID = userID
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		return
	}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		slot := &models.UnavailableSlot{
			UnavailableDate: d.Format("2006-01-02"),
			StartTime:       input.StartTime,
			EndTime:         input.EndTime,
			UID:             input.UID,
		}

		err := dao.CreateUnavailableSlots(slot)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Failed to create record for date: "+d.Format("2006-01-02"), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUnavailabilitySlots(w http.ResponseWriter, r *http.Request) {
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
	unavailabilitySlots, err := dao.GetUnavailableSlots(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("unavailability with user_id %d not found", userID) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve unavailability slots", http.StatusInternalServerError)
			log.Printf("failed to retrieve unavailability slots: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unavailabilitySlots)
}
