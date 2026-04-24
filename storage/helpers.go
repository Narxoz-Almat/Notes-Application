package storage

import (
	"sort"
	"strings"

	"notes-app/models"
)

func paginate(total, page, limit int) (int, int) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start > total {
		return total, total
	}
	end := start + limit
	if end > total {
		end = total
	}
	return start, end
}

func matchBook(book models.Book, filter BookFilter) bool {
	if filter.Category != "" && !strings.Contains(strings.ToLower(book.Category.Name), strings.ToLower(filter.Category)) {
		return false
	}
	if filter.Author != "" && !strings.Contains(strings.ToLower(book.Author.Name), strings.ToLower(filter.Author)) {
		return false
	}
	if filter.Title != "" && !strings.Contains(strings.ToLower(book.Title), strings.ToLower(filter.Title)) {
		return false
	}
	if filter.MinPrice != nil && book.Price < *filter.MinPrice {
		return false
	}
	if filter.MaxPrice != nil && book.Price > *filter.MaxPrice {
		return false
	}
	return true
}

func sortBooks(books []models.Book) {
	sort.Slice(books, func(i, j int) bool {
		return books[i].ID < books[j].ID
	})
}

func sortAuthors(authors []models.Author) {
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].ID < authors[j].ID
	})
}

func sortCategories(categories []models.Category) {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})
}
