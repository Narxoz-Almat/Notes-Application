package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"notes-app/models"
	"notes-app/storage"
)

type App struct {
	repo      storage.Repository
	jwtSecret string
}

func NewApp(repo storage.Repository, jwtSecret string) *App {
	return &App{repo: repo, jwtSecret: jwtSecret}
}

type bookRequest struct {
	Title      string  `json:"title" binding:"required"`
	AuthorID   uint    `json:"author_id" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
}

type authorRequest struct {
	Name string `json:"name" binding:"required"`
}

type categoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (a *App) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		parsedClaims, ok := token.Claims.(*claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("userID", parsedClaims.UserID)
		c.Next()
	}
}

func (a *App) userIDFromContext(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := value.(uint)
	return userID, ok
}

func (a *App) userIDFromAuthorization(c *gin.Context) (uint, bool) {
	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return 0, false
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}

	parsedClaims, ok := token.Claims.(*claims)
	if !ok {
		return 0, false
	}

	return parsedClaims.UserID, true
}

func (a *App) parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return uint(value), true
}

func (a *App) issueToken(user models.User) (string, error) {
	issuedAt := jwt.NewNumericDate(time.Now().UTC())
	expiresAt := jwt.NewNumericDate(time.Now().UTC().AddDate(0, 0, 7))
	claims := claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func repoErrorStatus(err error) (int, string) {
	switch err {
	case storage.ErrNotFound:
		return http.StatusNotFound, err.Error()
	case storage.ErrAlreadyExists:
		return http.StatusConflict, err.Error()
	case storage.ErrInvalidRelated:
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
