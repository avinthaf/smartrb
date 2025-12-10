package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"smartrb.com/db"
	"smartrb.com/flashcards"
	"smartrb.com/users"
)

func handleGetFlashcardDecks(c *gin.Context) {

	checkAuthorization(c)

	flashcardDecks, err := flashcards.GetFlashcardDecks(db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get flashcard decks",
			"details": err.Error(),
		})
		return
	}

	// Loop over flashcardDecks and get ids into an array
	// var flashcardDeckIds []string
	// for i := range flashcardDecks {
	// 	flashcardDeckIds = append(flashcardDeckIds, flashcardDecks[i].Id)
	// }

	// productsCategories, err := categories.GetProductsCategoriesByProductIds(flashcardDeckIds, db.Default())
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Failed to get products categories",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, flashcardDecks)

}

func handleGetFlashcardsByDeckId(c *gin.Context) {

	checkAuthorization(c)

	deckId := c.Param("deckId")

	fmt.Println(deckId)

	flashcards, err := flashcards.GetFlashcardsByDeckId(db.Default(), deckId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get flashcards",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, flashcards)

}

func handleCreateFlashcardDeckSession(c *gin.Context) {
	// Check authorization and get user ID
	userId, ok := checkAuthorization(c)
	if !ok {
		// checkAuthorization already sent the error response
		return
	}

	// Validate user ID is not empty
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user by external id
	user, err := users.GetUserByExternalId(userId, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user by external id",
			"details": err.Error(),
		})
		return
	}

	// Validate user ID from database
	if user.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	deckId := c.Param("deckId")
	sessionId := c.Param("sessionId")

	flashcardDeckSession, err := flashcards.CreateFlashcardDeckSession(db.Default(), sessionId, deckId, user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create flashcard deck session",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard deck session created successfully",
		"session": flashcardDeckSession,
	})
}

func handleCreateFlashcardScore(c *gin.Context) {
	// Check authorization and get user ID
	userId, ok := checkAuthorization(c)
	if !ok {
		// checkAuthorization already sent the error response
		return
	}

	// Validate user ID is not empty
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user by external id
	user, err := users.GetUserByExternalId(userId, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user by external id",
			"details": err.Error(),
		})
		return
	}

	// Validate user ID from database
	if user.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var flashcardScoreReq flashcards.FlashcardScoreRequest

	if err := c.ShouldBindJSON(&flashcardScoreReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid flashcard score",
			"details": err.Error(),
		})
		return
	}

	// Validate card ID is not empty
	if flashcardScoreReq.CardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Card ID is required",
		})
		return
	}

	// Create the flashcard score
	if err := flashcards.CreateFlashcardScore(db.Default(), user.Id, flashcardScoreReq.CardId, flashcardScoreReq.Score, flashcardScoreReq.SessionId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create flashcard score",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard score created successfully",
	})

}
