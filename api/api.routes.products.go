package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"smartrb.com/db"
	"smartrb.com/products"
	"smartrb.com/users"
)

func handleCreateProductRating(c *gin.Context) {

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

	var req products.CreateProductRatingRequest

	// Bind JSON to struct
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	fmt.Println("User ID from database:", user.Id)
	fmt.Println("Product ID from request:", req.ProductId)
	fmt.Println("Rating from request:", req.Rating)

	// Create product rating
	err = products.CreateProductRating(db.Default(), req.ProductId, user.Id, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create product rating",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product rating created successfully",
	})
}
