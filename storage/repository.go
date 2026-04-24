package storage

import "notes-app/models"

type BookFilter struct {
	Category string
	Author   string
	Title    string
	MinPrice *float64
	MaxPrice *float64
}

type Repository interface {
	CreateBook(book *models.Book) error
	ListBooks(filter BookFilter, page, limit int) ([]models.Book, int64, error)
	GetBook(id uint) (*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id uint) error

	ListAuthors() ([]models.Author, error)
	CreateAuthor(author *models.Author) error

	ListCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error

	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)

	AddFavorite(userID, bookID uint) error
	RemoveFavorite(userID, bookID uint) error
	ListFavoriteBooks(userID uint, page, limit int) ([]models.Book, int64, error)
}
