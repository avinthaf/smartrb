package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smartrb.com/db"
	"smartrb.com/interests"
)

func handleGetInterestsByUserId(c *gin.Context) {

	userId, _ := checkAuthorization(c)

	interests, err := interests.GetInterestsByUserId(userId, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get interests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, interests)

}


func handleCreateInterests(c *gin.Context) {

	checkAuthorization(c)

	var reqs []interests.CreateInterestRequest

	// Bind JSON to struct
	err := c.ShouldBindJSON(&reqs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create interest
	err = interests.CreateInterests(reqs, db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create interest",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Interests created successfully",
	})

}
