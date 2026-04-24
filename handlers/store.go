package handlers

import (
	"sync"

	"notes-app/models"
)

var (
	notes = []models.Note{}
	users = []models.User{}
	tags  = []models.Tag{}

	noteID = 1
	userID = 1
	tagID  = 1

	storeMu sync.RWMutex
)
