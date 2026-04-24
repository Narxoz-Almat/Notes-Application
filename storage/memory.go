package storage

import (
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"notes-app/models"
)

type MemoryRepository struct {
	mu sync.RWMutex

	books      map[uint]models.Book
	authors    map[uint]models.Author
	categories map[uint]models.Category
	users      map[uint]models.User
	favorites  map[uint]map[uint]time.Time

	nextBookID     uint
	nextAuthorID   uint
	nextCategoryID uint
	nextUserID     uint
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		books:          make(map[uint]models.Book),
		authors:        make(map[uint]models.Author),
		categories:     make(map[uint]models.Category),
		users:          make(map[uint]models.User),
		favorites:      make(map[uint]map[uint]time.Time),
		nextBookID:     1,
		nextAuthorID:   1,
		nextCategoryID: 1,
		nextUserID:     1,
	}
}

func (r *MemoryRepository) CreateAuthor(author *models.Author) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, current := range r.authors {
		if strings.EqualFold(current.Name, author.Name) {
			return ErrAlreadyExists
		}
	}
	author.ID = r.nextAuthorID
	r.nextAuthorID++
	author.CreatedAt = time.Now().UTC()
	author.UpdatedAt = author.CreatedAt
	r.authors[author.ID] = *author
	return nil
}

func (r *MemoryRepository) ListAuthors() ([]models.Author, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	authors := make([]models.Author, 0, len(r.authors))
	for _, author := range r.authors {
		authors = append(authors, author)
	}
	sortAuthors(authors)
	return authors, nil
}

func (r *MemoryRepository) CreateCategory(category *models.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, current := range r.categories {
		if strings.EqualFold(current.Name, category.Name) {
			return ErrAlreadyExists
		}
	}
	category.ID = r.nextCategoryID
	r.nextCategoryID++
	category.CreatedAt = time.Now().UTC()
	category.UpdatedAt = category.CreatedAt
	r.categories[category.ID] = *category
	return nil
}

func (r *MemoryRepository) ListCategories() ([]models.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	categories := make([]models.Category, 0, len(r.categories))
	for _, category := range r.categories {
		categories = append(categories, category)
	}
	sortCategories(categories)
	return categories, nil
}

func (r *MemoryRepository) CreateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, current := range r.users {
		if strings.EqualFold(current.Email, user.Email) {
			return ErrAlreadyExists
		}
	}
	user.ID = r.nextUserID
	r.nextUserID++
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.users[user.ID] = *user
	return nil
}

func (r *MemoryRepository) FindUserByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			copy := user
			return &copy, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryRepository) FindUserByID(id uint) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := user
	return &copy, nil
}

func (r *MemoryRepository) CreateBook(book *models.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.authors[book.AuthorID]; !ok {
		return ErrInvalidRelated
	}
	if _, ok := r.categories[book.CategoryID]; !ok {
		return ErrInvalidRelated
	}
	book.ID = r.nextBookID
	r.nextBookID++
	book.CreatedAt = time.Now().UTC()
	book.UpdatedAt = book.CreatedAt
	r.books[book.ID] = *book
	return nil
}

func (r *MemoryRepository) ListBooks(filter BookFilter, page, limit int) ([]models.Book, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]models.Book, 0, len(r.books))
	for _, book := range r.books {
		expanded, ok := r.expandBookLocked(book)
		if !ok || !matchBook(expanded, filter) {
			continue
		}
		books = append(books, expanded)
	}
	sortBooks(books)
	total := int64(len(books))
	start, end := paginate(len(books), page, limit)
	return books[start:end], total, nil
}

func (r *MemoryRepository) GetBook(id uint) (*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, ok := r.books[id]
	if !ok {
		return nil, ErrNotFound
	}
	expanded, ok := r.expandBookLocked(book)
	if !ok {
		return nil, ErrInvalidRelated
	}
	return &expanded, nil
}

func (r *MemoryRepository) UpdateBook(book *models.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.books[book.ID]
	if !ok {
		return ErrNotFound
	}
	if _, ok := r.authors[book.AuthorID]; !ok {
		return ErrInvalidRelated
	}
	if _, ok := r.categories[book.CategoryID]; !ok {
		return ErrInvalidRelated
	}
	book.CreatedAt = current.CreatedAt
	book.UpdatedAt = time.Now().UTC()
	r.books[book.ID] = *book
	return nil
}

func (r *MemoryRepository) DeleteBook(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return ErrNotFound
	}
	delete(r.books, id)
	for userID := range r.favorites {
		delete(r.favorites[userID], id)
	}
	return nil
}

func (r *MemoryRepository) AddFavorite(userID, bookID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[userID]; !ok {
		return ErrNotFound
	}
	if _, ok := r.books[bookID]; !ok {
		return ErrNotFound
	}
	if r.favorites[userID] == nil {
		r.favorites[userID] = make(map[uint]time.Time)
	}
	if _, exists := r.favorites[userID][bookID]; exists {
		return ErrAlreadyExists
	}
	r.favorites[userID][bookID] = time.Now().UTC()
	return nil
}

func (r *MemoryRepository) RemoveFavorite(userID, bookID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.favorites[userID]; !ok {
		return ErrNotFound
	}
	if _, ok := r.favorites[userID][bookID]; !ok {
		return ErrNotFound
	}
	delete(r.favorites[userID], bookID)
	return nil
}

func (r *MemoryRepository) ListFavoriteBooks(userID uint, page, limit int) ([]models.Book, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	bookIDs, ok := r.favorites[userID]
	if !ok {
		return []models.Book{}, 0, nil
	}
	type favoriteEntry struct {
		bookID    uint
		createdAt time.Time
	}
	entries := make([]favoriteEntry, 0, len(bookIDs))
	for bookID, createdAt := range bookIDs {
		entries = append(entries, favoriteEntry{bookID: bookID, createdAt: createdAt})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].createdAt.After(entries[j].createdAt)
	})
	total := int64(len(entries))
	start, end := paginate(len(entries), page, limit)
	books := make([]models.Book, 0, end-start)
	for _, entry := range entries[start:end] {
		book, ok := r.books[entry.bookID]
		if !ok {
			continue
		}
		expanded, ok := r.expandBookLocked(book)
		if !ok {
			continue
		}
		books = append(books, expanded)
	}
	return books, total, nil
}

func (r *MemoryRepository) expandBookLocked(book models.Book) (models.Book, bool) {
	author, ok := r.authors[book.AuthorID]
	if !ok {
		return models.Book{}, false
	}
	category, ok := r.categories[book.CategoryID]
	if !ok {
		return models.Book{}, false
	}
	book.Author = author
	book.Category = category
	return book, true
}

func (r *MemoryRepository) verifyPassword(user *models.User, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plain))
}
