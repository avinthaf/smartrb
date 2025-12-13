package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"smartrb.com/categories"
	"smartrb.com/db"
	"smartrb.com/flashcards"
	"smartrb.com/users"
)

type FlashcardDeckWithCategory struct {
	Id            string                 `json:"id"`
	Title         string                 `json:"title"`
	Description   *string                `json:"description,omitempty"`
	UserId        *string                `json:"user_id,omitempty"`
	PublishStatus string                 `json:"publish_status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Categories    []*categories.Category `json:"categories,omitempty"`
}

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
	var flashcardDeckIds []string
	for i := range flashcardDecks {
		flashcardDeckIds = append(flashcardDeckIds, flashcardDecks[i].Id)
	}

	productCategories, err := categories.GetProductCategoriesByProductIds(flashcardDeckIds, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get products categories",
			"details": err.Error(),
		})
		return
	}

	// Loop over productCategories and get category ids into an array
	var categoryIds []string
	for i := range productCategories {
		categoryIds = append(categoryIds, productCategories[i].CategoryId)
	}

	categoryList, err := categories.GetCategoriesByIds(categoryIds, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get categories",
			"details": err.Error(),
		})
		return
	}

	// Create a new flashcard deck struct with category
	// Initialize the slice with the correct length
	flashcardDecksWithCategory := make([]FlashcardDeckWithCategory, len(flashcardDecks))

	// Match categories to flashcardDecks
	// productCategories product id maps to flashcardDeck id and category id maps to category id
	for i := range flashcardDecks {
		// Copy the basic flashcard deck data
		flashcardDecksWithCategory[i] = FlashcardDeckWithCategory{
			Id:            flashcardDecks[i].Id,
			Title:         flashcardDecks[i].Title,
			Description:   flashcardDecks[i].Description,
			UserId:        flashcardDecks[i].UserId,
			PublishStatus: flashcardDecks[i].PublishStatusId,
			CreatedAt:     flashcardDecks[i].CreatedAt,
			UpdatedAt:     flashcardDecks[i].UpdatedAt,
			Categories:    []*categories.Category{}, // Initialize empty pointer slice
		}

		// Find and assign all matching categories
		for j := range productCategories {
			if flashcardDecks[i].Id == productCategories[j].ProductId {
				for k := range categoryList {
					if productCategories[j].CategoryId == categoryList[k].Id {
						// Append the category pointer to the slice (don't break)
						flashcardDecksWithCategory[i].Categories = append(flashcardDecksWithCategory[i].Categories, &categoryList[k])
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, flashcardDecksWithCategory)

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

func handleGetFlashcardDeckSessionsByUserId(c *gin.Context) {
	// Check authorization and get user ID
	externalUserId, ok := checkAuthorization(c)
	if !ok {
		// checkAuthorization already sent the error response
		return
	}

	// Validate user ID is not empty
	if externalUserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user by external id
	user, err := users.GetUserByExternalId(externalUserId, db.Default())
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

	flashcardSessions, err := flashcards.GetFlashcardDeckSessionsByUserId(db.Default(), user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get flashcard sessions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, flashcardSessions)
}

func handleGetFlashcardScoresBySessionId(c *gin.Context) {
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

	sessionId := c.Param("sessionId")

	flashcardScores, err := flashcards.GetFlashcardScoresBySessionId(db.Default(), sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get flashcard scores",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, flashcardScores)
}

func handleCreateFlashcardDeck(c *gin.Context) {
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

	// Parse request body
	var deckReq flashcards.CreateFlashcardDeckRequest

	if err := c.ShouldBindJSON(&deckReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid flashcard deck request",
			"details": err.Error(),
		})
		return
	}

	flashcardDeck, err := flashcards.CreateFlashcardDeck(db.Default(), user.Id, deckReq.Title, deckReq.Description, deckReq.PublishStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create flashcard deck",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, flashcardDeck)

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

	// Parse request body
	var sessionReq flashcards.FlashcardDeckSessionRequest

	if err := c.ShouldBindJSON(&sessionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid flashcard deck session request",
			"details": err.Error(),
		})
		return
	}

	// Validate deck ID is not empty
	if sessionReq.DeckId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Deck ID is required",
		})
		return
	}

	// Validate session ID is not empty
	if sessionReq.SessionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session ID is required",
		})
		return
	}

	flashcardDeckSession, err := flashcards.CreateFlashcardDeckSession(db.Default(), sessionReq.SessionId, sessionReq.DeckId, user.Id)

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
