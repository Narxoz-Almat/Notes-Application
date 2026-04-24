package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"notes-app/models"

	"github.com/gorilla/mux"
)

// GET /notes
func GetNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storeMu.RLock()
	defer storeMu.RUnlock()

	_ = json.NewEncoder(w).Encode(notes)
}

// POST /notes
func CreateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	note.ID = noteID
	note.CreatedAt = time.Now()
	noteID++

	notes = append(notes, note)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(note)
}

// GET /notes/{id}
func GetNoteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	storeMu.RLock()
	defer storeMu.RUnlock()

	for _, note := range notes {
		if note.ID == id {
			_ = json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

// PUT /notes/{id}
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	var payload models.Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	for i, note := range notes {
		if note.ID == id {
			payload.ID = id
			payload.CreatedAt = note.CreatedAt
			notes[i] = payload

			_ = json.NewEncoder(w).Encode(payload)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

// DELETE /notes/{id}
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}
