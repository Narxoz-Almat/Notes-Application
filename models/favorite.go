package models

import "time"

type FavoriteBook struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey;column:user_id"`
	BookID    uint      `json:"book_id" gorm:"primaryKey;column:book_id"`
	CreatedAt time.Time `json:"created_at"`
	Book      Book      `json:"book,omitempty" gorm:"foreignKey:BookID"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}
