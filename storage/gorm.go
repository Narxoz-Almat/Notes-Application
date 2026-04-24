package storage

import (
	"time"

	"gorm.io/gorm"

	"notes-app/models"
)

type GormRepository struct {
	db *gorm.DB
}

func (r *GormRepository) CreateAuthor(author *models.Author) error {
	var existing models.Author
	if err := r.db.Where("LOWER(name) = LOWER(?)", author.Name).First(&existing).Error; err == nil {
		return ErrAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(author).Error
}

func (r *GormRepository) ListAuthors() ([]models.Author, error) {
	var authors []models.Author
	if err := r.db.Order("id ASC").Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *GormRepository) CreateCategory(category *models.Category) error {
	var existing models.Category
	if err := r.db.Where("LOWER(name) = LOWER(?)", category.Name).First(&existing).Error; err == nil {
		return ErrAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(category).Error
}

func (r *GormRepository) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *GormRepository) CreateUser(user *models.User) error {
	var existing models.User
	if err := r.db.Where("LOWER(email) = LOWER(?)", user.Email).First(&existing).Error; err == nil {
		return ErrAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(user).Error
}

func (r *GormRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) CreateBook(book *models.Book) error {
	if err := r.ensureAuthorAndCategory(book.AuthorID, book.CategoryID); err != nil {
		return err
	}
	return r.db.Create(book).Error
}

func (r *GormRepository) ListBooks(filter BookFilter, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	if err := r.db.Preload("Author").Preload("Category").Order("id ASC").Find(&books).Error; err != nil {
		return nil, 0, err
	}
	filtered := make([]models.Book, 0, len(books))
	for _, book := range books {
		if matchBook(book, filter) {
			filtered = append(filtered, book)
		}
	}
	sortBooks(filtered)
	total := int64(len(filtered))
	start, end := paginate(len(filtered), page, limit)
	return filtered[start:end], total, nil
}

func (r *GormRepository) GetBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := r.db.Preload("Author").Preload("Category").First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &book, nil
}

func (r *GormRepository) UpdateBook(book *models.Book) error {
	if _, err := r.GetBook(book.ID); err != nil {
		return err
	}
	if err := r.ensureAuthorAndCategory(book.AuthorID, book.CategoryID); err != nil {
		return err
	}
	book.UpdatedAt = time.Now().UTC()
	return r.db.Save(book).Error
}

func (r *GormRepository) DeleteBook(id uint) error {
	result := r.db.Delete(&models.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepository) AddFavorite(userID, bookID uint) error {
	if _, err := r.FindUserByID(userID); err != nil {
		return err
	}
	if _, err := r.GetBook(bookID); err != nil {
		return err
	}
	var existing models.FavoriteBook
	if err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing).Error; err == nil {
		return ErrAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	favorite := models.FavoriteBook{UserID: userID, BookID: bookID, CreatedAt: time.Now().UTC()}
	return r.db.Create(&favorite).Error
}

func (r *GormRepository) RemoveFavorite(userID, bookID uint) error {
	result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.FavoriteBook{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepository) ListFavoriteBooks(userID uint, page, limit int) ([]models.Book, int64, error) {
	var favorites []models.FavoriteBook
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, 0, err
	}
	total := int64(len(favorites))
	start, end := paginate(len(favorites), page, limit)
	books := make([]models.Book, 0, end-start)
	for _, favorite := range favorites[start:end] {
		book, err := r.GetBook(favorite.BookID)
		if err != nil {
			continue
		}
		books = append(books, *book)
	}
	return books, total, nil
}

func (r *GormRepository) ensureAuthorAndCategory(authorID, categoryID uint) error {
	if _, err := r.findAuthorByID(authorID); err != nil {
		return err
	}
	if _, err := r.findCategoryByID(categoryID); err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) findAuthorByID(id uint) (*models.Author, error) {
	var author models.Author
	if err := r.db.First(&author, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidRelated
		}
		return nil, err
	}
	return &author, nil
}

func (r *GormRepository) findCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidRelated
		}
		return nil, err
	}
	return &category, nil
}
