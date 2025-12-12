package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"smartrb.com/db"
	"smartrb.com/fill_in_blanks"
	"smartrb.com/users"
	"smartrb.com/categories"
)

type FillInBlankDeckWithCategory struct {
	Id            string                   `json:"id"`
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	UserId        *string                  `json:"user_id,omitempty"`
	PublishStatus string                   `json:"publish_status"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Categories    []*categories.Category    `json:"categories,omitempty"`
}

func handleGetFillInBlankDecks(c *gin.Context) {

	checkAuthorization(c)

	fillInBlankDecks, err := fill_in_blanks.GetFillInBlankDecks(db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get fill in blank decks",
			"details": err.Error(),
		})
		return
	}

	// Loop over fillInBlankDecks and get ids into an array
	var fillInBlankDeckIds []string
	for i := range fillInBlankDecks {
		fillInBlankDeckIds = append(fillInBlankDeckIds, fillInBlankDecks[i].Id)
	}

	productCategories, err := categories.GetProductCategoriesByProductIds(fillInBlankDeckIds, db.Default())
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

	// Create a new fill in blank deck struct with category
	var fillInBlankDecksWithCategory []FillInBlankDeckWithCategory

	// Initialize the slice with the correct length
	fillInBlankDecksWithCategory = make([]FillInBlankDeckWithCategory, len(fillInBlankDecks))

	// Match categories to fillInBlankDecks
	// productCategories product id maps to fillInBlankDeck id and category id maps to category id
	for i := range fillInBlankDecks {
		// Copy the basic fill in blank deck data
		fillInBlankDecksWithCategory[i] = FillInBlankDeckWithCategory{
			Id:            fillInBlankDecks[i].Id,
			Title:         fillInBlankDecks[i].Title,
			Description:   fillInBlankDecks[i].Description,
			UserId:        fillInBlankDecks[i].UserId,
			PublishStatus: fillInBlankDecks[i].PublishStatus,
			CreatedAt:     fillInBlankDecks[i].CreatedAt,
			UpdatedAt:     fillInBlankDecks[i].UpdatedAt,
			Categories:    []*categories.Category{}, // Initialize empty pointer slice
		}

		// Find and assign all matching categories
		for j := range productCategories {
			if fillInBlankDecks[i].Id == productCategories[j].ProductId {
				for k := range categoryList {
					if productCategories[j].CategoryId == categoryList[k].Id {
						// Append the category pointer to the slice (don't break)
						fillInBlankDecksWithCategory[i].Categories = append(fillInBlankDecksWithCategory[i].Categories, &categoryList[k])
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, fillInBlankDecksWithCategory)

}

func handleGetFillInBlanksByDeckId(c *gin.Context) {

	checkAuthorization(c)

	deckId := c.Param("deckId")

	fmt.Println(deckId)

	fillInBlanks, err := fill_in_blanks.GetFillInBlanksByDeckId(db.Default(), deckId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get fill in blanks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fillInBlanks)

}

func handleGetFillInBlankDeckSessionsByUserId(c *gin.Context) {
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

	fillInBlankSessions, err := fill_in_blanks.GetFillInBlankDeckSessionsByUserId(db.Default(), user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get fill in blank sessions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fillInBlankSessions)
}

func handleGetFillInBlankScoresBySessionId(c *gin.Context) {
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

	fillInBlankScores, err := fill_in_blanks.GetFillInBlankScoresBySessionId(db.Default(), sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get fill in blank scores",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fillInBlankScores)
}

func handleCreateFillInBlankDeckSession(c *gin.Context) {
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
	var sessionReq fill_in_blanks.FillInBlankDeckSessionRequest

	if err := c.ShouldBindJSON(&sessionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fill in blank deck session request",
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

	fillInBlankDeckSession, err := fill_in_blanks.CreateFillInBlankDeckSession(db.Default(), sessionReq, user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create fill in blank deck session",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fill in blank deck session created successfully",
		"session": fillInBlankDeckSession,
	})
}

func handleCreateFillInBlankScore(c *gin.Context) {
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

	var fillInBlankScoreReq fill_in_blanks.FillInBlankScoreRequest

	if err := c.ShouldBindJSON(&fillInBlankScoreReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fill in blank score",
			"details": err.Error(),
		})
		return
	}

	// Validate fill in blank ID is not empty
	if fillInBlankScoreReq.FillInBlankId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fill in blank ID is required",
		})
		return
	}

	// Create the fill in blank score
	if err := fill_in_blanks.CreateFillInBlankScore(db.Default(), fillInBlankScoreReq, user.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create fill in blank score",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fill in blank score created successfully",
	})

}