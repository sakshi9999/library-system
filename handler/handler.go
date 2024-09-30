package handler

import (
	"library-system/models"
	"library-system/stats"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	AddBookAPI    = "ADD_BOOK"
	ListBookAPI   = "LIST_BOOK_API"
	ReturnBookAPI = "RETURN_BOOK_API"
	BorrowBookAPI = "BORROW_BOOK_API"
)

type CommandHandler struct {
	Db *gorm.DB
}

func NewCommandHandler(db *gorm.DB) CommandHandler {
	return CommandHandler{Db: db}
}

func (d *CommandHandler) AddBook(c *gin.Context) {
	apiStartTime := time.Now()
	stats.AddBookApiCounter.Inc()
	defer stats.ApiElapsedTime.WithLabelValues(AddBookAPI).Set(float64(time.Since(apiStartTime).Milliseconds()))

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.AvailableCopies = book.Copies
	d.Db.Create(&book)
	c.JSON(http.StatusOK, book)
}

// ListBooks - GET /books
func (d *CommandHandler) ListBooks(c *gin.Context) {
	apiStartTime := time.Now()
	stats.GetBookApiCounter.Inc()
	defer stats.ApiElapsedTime.WithLabelValues(ListBookAPI).Set(float64(time.Since(apiStartTime).Milliseconds()))

	var books []models.Book
	d.Db.Find(&books)
	c.JSON(http.StatusOK, books)
}

// BorrowBook - POST /books/:id/borrow
func (d *CommandHandler) BorrowBook(c *gin.Context) {
	apiStartTime := time.Now()
	stats.BorrowBookApiCounter.Inc()
	defer stats.ApiElapsedTime.WithLabelValues(BorrowBookAPI).Set(float64(time.Since(apiStartTime).Milliseconds()))

	var book models.Book
	if err := d.Db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.AvailableCopies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No copies available"})
		return
	}

	var borrower models.Borrower
	if err := c.ShouldBindJSON(&borrower); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrower.BookID = book.ID
	borrower.Status = "borrowed"

	d.Db.Create(&borrower)
	book.AvailableCopies -= 1
	d.Db.Save(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Book borrowed successfully"})
}

// ReturnBook - POST /books/:id/return
func (d *CommandHandler) ReturnBook(c *gin.Context) {
	apiStartTime := time.Now()
	stats.ReturnBookApiCounter.Inc()
	defer stats.ApiElapsedTime.WithLabelValues(ReturnBookAPI).Set(float64(time.Since(apiStartTime).Milliseconds()))

	var book models.Book
	if err := d.Db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var borrower models.Borrower
	if err := c.ShouldBindJSON(&borrower); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := d.Db.Where("book_id = ? AND name = ? AND status = 'borrowed'", book.ID, borrower.Name).First(&borrower).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrow record not found"})
		return
	}

	borrower.Status = "returned"
	d.Db.Save(&borrower)

	book.AvailableCopies += 1
	d.Db.Save(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
