package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/Goldmite/project_shelf/internal/models/dto"
	"github.com/Goldmite/project_shelf/internal/services"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bs *services.BookService) *BookHandler {
	return &BookHandler{bookService: bs}
}

func (bookHandler *BookHandler) GetBookByIsbnHandler(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := bookHandler.bookService.GetBookByIsbn(isbn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("book with ISBN %s not found", isbn)})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bookHandler *BookHandler) AddNewBookForUserHandler(c *gin.Context) {
	var req dto.UserBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := bookHandler.bookService.AddNewBookForUser(req.UserId, req.Isbn)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate user book: " + err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (bookHandler *BookHandler) GetAllUserBooksHandler(c *gin.Context) {
	userId := c.Param("id")
	userBooks, err := bookHandler.bookService.GetAllUserBooks(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userBooks)
}

func (bookHandler *BookHandler) GetAllGroupBooksHandler(c *gin.Context) {
	groupId := c.Param("id")
	groupBooks, err := bookHandler.bookService.GetAllGroupBooks(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groupBooks)
}

func (bookHandler *BookHandler) GetRecentlyReadBookHandler(c *gin.Context) {
	userId := c.Param("id")
	recentBook, err := bookHandler.bookService.GetRecentlyReadBook(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No recently read book in the past month"})
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recentBook)
}
