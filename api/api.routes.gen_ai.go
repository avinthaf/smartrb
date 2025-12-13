package api

import (
	"net/http"
	
	"github.com/gin-gonic/gin"

	"smartrb.com/gen_ai"
	"smartrb.com/db"
)

func handleGenAIPrompt(c *gin.Context) {
	// Get prompt from request body
	var requestData struct {
		Prompt string `json:"prompt"`
	}
	
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// Call the AI service
	result, err := gen_ai.CreateAIContent(requestData.Prompt, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"result": result})
}
