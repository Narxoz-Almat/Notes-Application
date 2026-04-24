package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"notes-app/models"
	"notes-app/storage"
)

func (a *App) ListBooks(c *gin.Context) {
	page, limit := paginationFromContext(c)
	filter := storage.BookFilter{
		Category: c.Query("category"),
		Author:   c.Query("author"),
		Title:    c.Query("title"),
	}
	if value := c.Query("min_price"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			filter.MinPrice = &parsed
		}
	}
	if value := c.Query("max_price"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			filter.MaxPrice = &parsed
		}
	}

	books, total, err := a.repo.ListBooks(filter, page, limit)
	if err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (a *App) CreateBook(c *gin.Context) {
	var req bookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title:      req.Title,
		AuthorID:   req.AuthorID,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
	if err := a.repo.CreateBook(&book); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	created, err := a.repo.GetBook(book.ID)
	if err == nil {
		c.JSON(http.StatusCreated, created)
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (a *App) GetBook(c *gin.Context) {
	id, ok := a.parseUintParam(c, "id")
	if !ok {
		return
	}

	book, err := a.repo.GetBook(id)
	if err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (a *App) UpdateBook(c *gin.Context) {
	id, ok := a.parseUintParam(c, "id")
	if !ok {
		return
	}

	var req bookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		ID:         id,
		Title:      req.Title,
		AuthorID:   req.AuthorID,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
	if err := a.repo.UpdateBook(&book); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	updated, err := a.repo.GetBook(id)
	if err == nil {
		c.JSON(http.StatusOK, updated)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (a *App) DeleteBook(c *gin.Context) {
	id, ok := a.parseUintParam(c, "id")
	if !ok {
		return
	}

	if err := a.repo.DeleteBook(id); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *App) DispatchBookRoutes(c *gin.Context) {
	path := strings.TrimPrefix(c.Param("path"), "/")
	if path == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	parts := strings.Split(path, "/")
	switch {
	case len(parts) == 1 && parts[0] == "favorites" && c.Request.Method == http.MethodGet:
		a.ListFavoriteBooks(c)
	case len(parts) == 1:
		c.Params = append(c.Params, gin.Param{Key: "id", Value: parts[0]})
		switch c.Request.Method {
		case http.MethodGet:
			a.GetBook(c)
		case http.MethodPut:
			a.UpdateBook(c)
		case http.MethodDelete:
			a.DeleteBook(c)
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
		}
	case len(parts) == 2 && parts[1] == "favorites":
		c.Params = append(c.Params, gin.Param{Key: "bookId", Value: parts[0]})
		switch c.Request.Method {
		case http.MethodPut:
			a.AddFavorite(c)
		case http.MethodDelete:
			a.RemoveFavorite(c)
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
		}
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func paginationFromContext(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
