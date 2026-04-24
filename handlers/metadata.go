package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"notes-app/models"
)

func (a *App) ListAuthors(c *gin.Context) {
	authors, err := a.repo.ListAuthors()
	if err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": authors})
}

func (a *App) CreateAuthor(c *gin.Context) {
	var req authorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	author := models.Author{Name: req.Name}
	if err := a.repo.CreateAuthor(&author); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusCreated, author)
}

func (a *App) ListCategories(c *gin.Context) {
	categories, err := a.repo.ListCategories()
	if err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (a *App) CreateCategory(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category := models.Category{Name: req.Name}
	if err := a.repo.CreateCategory(&category); err != nil {
		status, message := repoErrorStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusCreated, category)
}
