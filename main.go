package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{
		ID:       "1",
		Author:   "John",
		Title:    "The Hobbit",
		Quantity: 1,
	},
	{
		ID:       "2",
		Author:   "John",
		Title:    "The Book",
		Quantity: 2,
	},
	{
		ID:       "3",
		Author:   "John",
		Title:    "The STORY",
		Quantity: 3,
	},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}
func getABookById(id string) (*book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found: " + id)
}
func getBookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getABookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, book)

}
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}
	book, err := getABookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if book.Quantity < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}
	book, err := getABookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)

}
func main() {
	var router = gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run(":8080")

}
