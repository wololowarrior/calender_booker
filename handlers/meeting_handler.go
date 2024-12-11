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

func CreateMeetings(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var meeting models.Meeting
	err := decoder.Decode(&meeting)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	_, err = dao.GetEvent(meeting.EventID)
	if err != nil {
		log.Print(err.Error())
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
	}
	err = dao.CreateMeeting(&meeting)
	if err != nil {
		log.Print(err.Error())
		if err.Error() == "unavailable during selected time" {
			http.Error(w, "User not available. Select another slot", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Failed to save meeting", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(meeting)
}

func GetMeetingsFromEvent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var meeting models.Meeting
	err := decoder.Decode(&meeting)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	eventID := meeting.EventID
	if eventID == 0 {
		http.Error(w, "Invalid event_id", http.StatusBadRequest)
		return
	}
	userID := meeting.UID
	if userID == 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
	}

	date := meeting.Date

	event, err := dao.GetEvent(eventID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		return
	}
	var slot string
	if event.Slots == nil {
		slot = "60"
	} else {
		slot = *event.Slots
	}
	slots, err := dao.GetSlottedMeetingsRecommendation(userID, slot, date)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

func UpdateMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meetingIDstr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'meetingID' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(meetingIDstr)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}

	m := models.Meeting{ID: id}

	err = dao.GetMeeting(&m)
	if err != nil {
		if err.Error() == "meeting not found" {
			http.Error(w, "Meeting not found", http.StatusNotFound)
		}
		http.Error(w, "Failed to get Meeting", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&m)
	log.Println(m)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	err = dao.UpdateMeeting(&m)
	if err != nil {
		if err.Error() == "invalid meeting" {
			http.Error(w, "Meeting is not valid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to save meeting", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func GetMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meetingIDstr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'meetingID' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(meetingIDstr)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}

	m := models.Meeting{ID: id}

	err = dao.GetMeeting(&m)
	if err != nil {
		if err.Error() == "meeting not found" {
			http.Error(w, "Meeting not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get Meeting", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

func DeleteMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meetingIDstr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'meetingID' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(meetingIDstr)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}

	m := models.Meeting{ID: id}

	err = dao.GetMeeting(&m)
	if err != nil {
		if err.Error() == "meeting not found" {
			http.Error(w, "Meeting not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get Meeting", http.StatusInternalServerError)
		return
	}
	err = dao.DeleteMeeting(&m)
	if err != nil {
		http.Error(w, "Failed to delete meeting", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetMeetingsForAUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'userID' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	date := queryParams.Get("date")
	meetings, err := dao.GetBookedMeetings(id, date)
	if err != nil {
		http.Error(w, "Failed to get Meetings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meetings)
}
