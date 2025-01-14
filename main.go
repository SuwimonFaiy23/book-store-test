package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	Author     string  `json:"author"`
	Year       string  `json:"year"`
	Amount     int64   `json:"amount"`
	Price      float64 `json:"price"`
	CategoryID int64   `json:"category_id"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// In-memory book storage
var books []Book
var nextID int64 = 1
var category = []Category{
	{
		ID:   1,
		Name: "Fiction",
	},
	{
		ID:   2,
		Name: "Non - Fiction",
	},
	{
		ID:   3,
		Name: "Reference Book",
	},
}

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Define the API routes
	// POST /books - Create a new book
	app.Post("/books", createBook)

	// PUT /books/:id - Update a book by ID
	app.Put("/books/:id", updateBook)

	// DELETE /books/:id - Delete a book by ID
	app.Delete("/books/:id", deleteBook)

	// GET /books/search - Search for books by title or author
	app.Get("/books/search", searchBooks)

	// Start the server
	err := app.Listen(":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// // Handler to create a new book
func createBook(c *fiber.Ctx) error {
	var newBook Book
	// Parse the JSON request body
	if err := c.BodyParser(&newBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Set the book's ID and increment the next ID
	newBook.ID = nextID
	nextID++

	books = append(books, newBook)

	// Respond with the created book
	return c.Status(fiber.StatusCreated).JSON(newBook)
}

// Handler to update a book by ID
func updateBook(c *fiber.Ctx) error {
	var updateBook Book
	// Get book ID form path
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}
	// Parse the JSON request body for update data book
	if err := c.BodyParser(&updateBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Loop update book detail by id
	for _, v := range books {
		if v.ID == id {
			v.Title = updateBook.Title
			v.Author = updateBook.Author
			v.Price = updateBook.Price
			v.Year = updateBook.Year
			v.Amount = updateBook.Amount
			v.CategoryID = updateBook.CategoryID
			books = append(books, v)
			return c.Status(fiber.StatusOK).JSON(v)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
}

// Handler to delete a book by ID
func deleteBook(c *fiber.Ctx) error {
	// Get book ID form path
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	// Loop delete book by id
	for index, v := range books {
		if v.ID == id {
			books = append(books[:index], books[:index+1]...)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Delete book success"})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
}

// Handler to search books by title, author or category
func searchBooks(c *fiber.Ctx) error {
	// Extract query parameters
	title := c.Query("title")                // Default empty if not provided
	author := c.Query("author")              // Default empty if not provided
	categoryName := c.Query("category_name") // Default empty if not provided

	var result []Book
	var categoryID int64
	// Loop find book by title, author or category
	for _, v := range books {
		if title != "" {
			if strings.Contains(strings.ToLower(v.Title), strings.ToLower(title)) {
				result = append(result, v)
			}
		} 
		 if author != "" {
			if strings.Contains(strings.ToLower(v.Author), strings.ToLower(author)) {
				result = append(result, v)
			}
		} 
		if categoryName != "" {
			for _, value := range category {
				if strings.Contains(strings.ToLower(value.Name), strings.ToLower(categoryName)) {
					categoryID = value.ID
				}
			}
			if categoryID == 0 {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No category found matching the search"})
			}

			if v.CategoryID == categoryID {
				result = append(result, v)
			}

		}
		// Check the number of books 
		if len(result) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No books found matching the search"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
