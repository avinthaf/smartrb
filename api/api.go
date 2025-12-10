package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
	"smartrb.com/auth"
)

// TODO: add middleware to check user permissions

func StartAPI() {

	// SETUP

	// Create auth client and JWKS provider
	auth_client := newAuthClient()

	// Start message receiver in background
	createMQSubscriber()

	// Create router
	router := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("CLIENT_URL")}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	// Load HTML templates
	// router.LoadHTMLGlob("courses/templates/*.html")

	// ROUTES

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Smartr API is running"})
	})

	// webhooks
	webhooks := router.Group("/webhooks")
	// Only accepts "authentication.email_verification_succeeded" event
	webhooks.POST("/signup", handleSignUp)
	// webhooks.POST("/login", handleLogin)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(auth.AuthMiddleware(auth_client))

	{
		// Auth routes
		// protected.GET("/login/checks", handlePostLoginChecks)

		// Categories Routes
		protected.GET("/categories/primary", handleGetPrimaryCategories)

		// Interests Routes
		protected.GET("/interests", handleGetInterestsByUserId)
		protected.POST("/interests", handleCreateInterests)

		// Flashcard Routes
		protected.GET("/flashcard_decks", handleGetFlashcardDecks)
		protected.GET("/flashcard_decks/:deckId", handleGetFlashcardsByDeckId)
		protected.POST("/flashcard_decks/:deckId/sessions/:sessionId", handleCreateFlashcardDeckSession)
		protected.POST("/flashcard_scores", handleCreateFlashcardScore)

		// // Gen UI
		// protected.POST("/gen_ui/courses/activities", handleCourseActvitiyGenUIPrompt)
	}

	router.Run(":8080")

}
