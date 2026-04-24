package handlers

import (
	"encoding/json"
	"net/http"

	"notes-app/models"
)

// GET /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storeMu.RLock()
	defer storeMu.RUnlock()

	_ = json.NewEncoder(w).Encode(users)
}

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	user.ID = uint(userID)
	userID++

	users = append(users, user)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
