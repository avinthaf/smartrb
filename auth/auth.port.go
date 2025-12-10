package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignUpMessage struct {
	Email     string    `json:"email"`
	UserId    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type LoginMessage struct {
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

type SignUpRequest struct {
	Email string `json:"email" binding:"required"`
}

type SignUpResult struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Success bool   `json:"success"`
}

type LoginResult struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
}


func GetAuthTopicKeys() []string {
	return auth_topic_keys
}

func SignUp(email string, user_id string, mq_callback MqCallback) SignUpResult {
	// Call SignUp service
	res := signUpService(email, user_id, mq_callback)

	return res
}

func Login(email string, mq_callback MqCallback) LoginResult {
	// Call Login service
	res := loginService(email, mq_callback)

	return res
}

func CreateJWTToken(id string) (string, error) {
	var secretKey = []byte("your-secret-key-here")

	claims := jwt.RegisteredClaims{
		Subject:   id, // Using ID as subject
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "your-app-name",
	}

	_ = claims

	// For custom claims with RegisteredClaims, create a custom struct or use MapClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(secretKey)
}

// Auth middleware for Supabase
func AuthMiddleware(auth_client *AuthClient) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with 'Bearer '"})
			c.Abort()
			return
		}

		// First, parse the token without verification to check the algorithm
		unverifiedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format: " + err.Error()})
			c.Abort()
			return
		}

		// Check the algorithm and use appropriate verification method
		if alg, ok := unverifiedToken.Header["alg"].(string); ok {
			if alg == "HS256" {
				// Use shared secret verification for legacy JWT secret
				err := verifyTokenWithSecret(tokenString, auth_client, c)
				if err != nil {
					return
				}
			} else if alg == "RS256" {
				// Use JWKS verification for asymmetric keys
				err := verifyTokenWithJWKS(tokenString, auth_client, c)
				if err != nil {
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unsupported signing algorithm: " + alg})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing algorithm"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func verifyTokenWithSecret(tokenString string, auth_client *AuthClient, c *gin.Context) error {
	// For legacy JWT secret, use HS256
	secretKey := []byte(auth_client.JWTSecret)
	if auth_client.JWTSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured for legacy tokens"})
		c.Abort()
		return fmt.Errorf("JWT secret not configured")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		c.Abort()
		return err
	}

	return validateClaimsAndSetContext(parsedToken, auth_client, c)
}

func verifyTokenWithJWKS(tokenString string, auth_client *AuthClient, c *gin.Context) error {
	// Extract key ID from token header and fetch public key from Supabase JWKS
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm (Supabase uses RS256 for asymmetric keys)
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get key ID from header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("token header missing key ID")
		}

		// Fetch the public key from Supabase JWKS endpoint
		return getAuthServicePublicKey(auth_client, kid)
	}

	parsedToken, err := jwt.Parse(tokenString, keyFunc)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		c.Abort()
		return err
	}

	return validateClaimsAndSetContext(parsedToken, auth_client, c)
}

func validateClaimsAndSetContext(parsedToken *jwt.Token, auth_client *AuthClient, c *gin.Context) error {
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Validate required Supabase claims
		expectedIssuer := auth_client.AuthServiceURL + "/auth/v1"
		if auth_client.AuthServiceURL == "" {
			expectedIssuer = "https://ebzxpsjqgupumzlmtygs.supabase.co/auth/v1" // fallback to hardcoded value
		}
		
		if iss, exists := claims["iss"]; !exists || iss != expectedIssuer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
			c.Abort()
			return fmt.Errorf("invalid issuer")
		}

		if aud, exists := claims["aud"]; !exists || aud != "authenticated" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience"})
			c.Abort()
			return fmt.Errorf("invalid audience")
		}

		// Extract user ID from token claims
		if user_id, exists := claims["sub"]; exists {
			c.Set("user_id", user_id)
		}

		// Extract email from user metadata if available
		if userMetadata, exists := claims["user_metadata"]; exists {
			if metadata, ok := userMetadata.(map[string]interface{}); ok {
				if email, exists := metadata["email"]; exists {
					c.Set("email", email)
				}
			}
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return fmt.Errorf("invalid claims")
	}

	return nil
}
