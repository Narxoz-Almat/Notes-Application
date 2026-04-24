package handlers

import (
	"encoding/json"
	"net/http"

	"notes-app/models"
)

// GET /tags
func GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storeMu.RLock()
	defer storeMu.RUnlock()

	_ = json.NewEncoder(w).Encode(tags)
}

// POST /tags
func CreateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	tag.ID = tagID
	tagID++

	tags = append(tags, tag)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(tag)
}
