package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smartrb.com/db"
	"smartrb.com/users"
)

func handleUpdateUser(c *gin.Context) {
	
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
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user by external id",
				"details": err.Error(),
			})
		}
		return
	}

	// Validate user ID from database
	if user.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get user info from request
	var req users.UpdateUserRequest

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Update user
	user, err = users.UpdateUser(user.Id, req.FirstName, req.LastName, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	// Return updated user
	c.JSON(http.StatusOK, user)
}