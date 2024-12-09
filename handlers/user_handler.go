package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"calendly_adventures/dao"
	"calendly_adventures/models"
	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing 'id' in URL path", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}
	user, err := dao.Get(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Error retrieving user: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
	}
	_, err = dao.Create(&user)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		if err.Error() == fmt.Sprintf("user with ID %d already exists", user.ID) {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
