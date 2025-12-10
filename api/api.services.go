package api

import (
	"os"
	"smartrb.com/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func newAuthClient() *auth.AuthClient {

	supabase_url := os.Getenv("SUPABASE_URL")
	jwt_secret := os.Getenv("SUPABASE_JWT_SECRET")

	return &auth.AuthClient{
		AuthServiceURL: supabase_url,
		JWTSecret:   jwt_secret, // Used for legacy HS256 tokens
	}
}


func checkAuthorization(c *gin.Context) (string, bool) {
	userId, userIdExists := c.Get("user_id")
	if !userIdExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return "", false
	}

	return userId.(string), true
}