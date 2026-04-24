package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) ListFavoriteBooks(c *gin.Context) {
	userID, ok := a.userIDFromAuthorization(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, limit := paginationFromContext(c)
	books, total, err := a.repo.ListFavoriteBooks(userID, page, limit)
	if err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books, "page": page, "limit": limit, "total": total})
}

func (a *App) AddFavorite(c *gin.Context) {
	userID, ok := a.userIDFromAuthorization(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	bookID, err := strconv.ParseUint(c.Param("bookId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}
	if err := a.repo.AddFavorite(userID, uint(bookID)); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *App) RemoveFavorite(c *gin.Context) {
	userID, ok := a.userIDFromAuthorization(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	bookID, err := strconv.ParseUint(c.Param("bookId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}
	if err := a.repo.RemoveFavorite(userID, uint(bookID)); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.Status(http.StatusNoContent)
}
