package api

import (
	"github.com/gin-gonic/gin"
	"smartrb.com/categories"
	"smartrb.com/db"
	"net/http"
)

func handleGetPrimaryCategories(c *gin.Context) {

	checkAuthorization(c)

	categories, err := categories.GetPrimaryCategories(db.Default())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
