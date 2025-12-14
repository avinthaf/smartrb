package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"smartrb.com/categories"
	"smartrb.com/db"
	"smartrb.com/gen_ai"
	"smartrb.com/flashcards"
	"smartrb.com/fill_in_blanks"
)

func handleFlashcardsGenAIPrompt(c *gin.Context) {

	checkAuthorization(c)

	// Get prompt from request body
	var requestData struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get all categories to provide context
	categories, err := categories.GetAllCategories(db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Prepare context for AI service
	categoriesContext := "Available categories: "
	for _, cat := range categories {
		categoriesContext += fmt.Sprintf("%s (%s) - ID: %s, ", cat.Name, cat.Description, cat.Id)
	}

	// Call the AI service with context
	result, err := gen_ai.CreateAIContent(flashcards.FlashcardSystemPrompt+ ", "+requestData.Prompt+ " Available Flashcard Categories: "+categoriesContext, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func handleFillInBlanksGenAIPrompt(c *gin.Context) {

	checkAuthorization(c)

	// Get prompt from request body
	var requestData struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get all categories to provide context
	categories, err := categories.GetAllCategories(db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Prepare context for AI service
	categoriesContext := "Available categories: "
	for _, cat := range categories {
		categoriesContext += fmt.Sprintf("%s (%s) - ID: %s, ", cat.Name, cat.Description, cat.Id)
	}

	// Call the AI service with context
	result, err := gen_ai.CreateAIContent(fill_in_blanks.FillInBlanksSystemPrompt+ ", "+requestData.Prompt+ " Available Fill-in-the-Blanks Categories: "+categoriesContext, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
