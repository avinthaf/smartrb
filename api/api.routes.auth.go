package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"smartrb.com/auth"
	"smartrb.com/mq"
)

func handleSignUp(c *gin.Context) {
	var webhook_event auth.AuthWebhookEvent
	if err := c.ShouldBindJSON(&webhook_event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use only the fields you need
	if webhook_event.Table == "users" && webhook_event.Type == "INSERT" {
		userID := webhook_event.Record.ID
		email := webhook_event.Record.Email
		fmt.Printf("New user: %s (%s)\n", email, userID)

		// Call your service
		auth.SignUp(email, userID, mq.Publish)
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}