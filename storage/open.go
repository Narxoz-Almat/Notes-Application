package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"notes-app/models"
)

func Open(dsn string) (Repository, error) {
	if dsn == "" {
		return NewMemoryRepository(), nil
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Author{}, &models.Category{}, &models.Book{}, &models.User{}, &models.FavoriteBook{}); err != nil {
		return nil, err
	}

	return &GormRepository{db: db}, nil
}
