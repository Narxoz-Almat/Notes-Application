package models

import "time"

type Book struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"not null"`
	AuthorID   uint      `json:"author_id" gorm:"not null;index"`
	CategoryID uint      `json:"category_id" gorm:"not null;index"`
	Price      float64   `json:"price" gorm:"not null"`
	Author     Author    `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Category   Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
