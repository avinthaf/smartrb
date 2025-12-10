package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
	Alg string   `json:"alg"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwksCache = struct {
	keys *JWKS
	exp  time.Time
}{}

type MqCallback func(topic string, routingKey string, message string)

var auth_topic_keys = []string{
	"auth.signup",
	"auth.login",
	"auth.logout",
}

func signUpService(email string, user_id string, mq_callback MqCallback) SignUpResult {

	// We just add the message to queue
	// TODO: implement any validations if needed

	// Send MQ message
	sendSignUpMessage(email, user_id, mq_callback)

	return SignUpResult{
		Id:      "",
		UserId:  user_id,
		Success: true,
	}
}

func loginService(email string, mq_callback MqCallback) LoginResult {

	// We just add the message to queue
	// TODO: implement any validations if needed

	// Send MQ message
	sendLoginMessage(email, mq_callback)

	return LoginResult{
		Id:      "",
		Success: true,
	}
}

func sendSignUpMessage(email string, user_id string, mq_callback MqCallback) {
	sign_up_message := SignUpMessage{
		Email:     email,
		UserId:    user_id,
		Timestamp: time.Now(),
	}

	json_data, err := json.Marshal(sign_up_message)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	mq_callback("events_topic", "auth.signup", string(json_data))
}

func sendLoginMessage(email string, mq_callback MqCallback) {
	login_message := LoginMessage{
		Email:     email,
		Timestamp: time.Now(),
	}

	json_data, err := json.Marshal(login_message)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	mq_callback("events_topic", "auth.login", string(json_data))
}


func getAuthServicePublicKey(auth_client *AuthClient, kid string) (*rsa.PublicKey, error) {
	// Check cache first
	if jwksCache.keys != nil && time.Now().Before(jwksCache.exp) {
		for _, key := range jwksCache.keys.Keys {
			if key.Kid == kid {
				return parseRSAPublicKey(&key)
			}
		}
	}

	// Construct Supabase JWKS URL
	authServiceURL := auth_client.AuthServiceURL
	if authServiceURL == "" {
		authServiceURL = os.Getenv("AUTH_SERVICE_URL")
	}
	
	jwksURL := authServiceURL + "/auth/v1/.well-known/jwks.json"

	// Fetch JWKS
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Cache the keys for 10 minutes (as recommended by Supabase docs)
	jwksCache.keys = &jwks
	jwksCache.exp = time.Now().Add(10 * time.Minute)

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(&key)
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

func parseRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	if len(jwk.X5c) > 0 {
		// Parse from X.509 certificate chain
		certPEM := "-----BEGIN CERTIFICATE-----\n" + jwk.X5c[0] + "\n-----END CERTIFICATE-----"
		return jwt.ParseRSAPublicKeyFromPEM([]byte(certPEM))
	}

	// Parse from n and e values
	return jwt.ParseRSAPublicKeyFromPEM([]byte(fmt.Sprintf(`-----BEGIN RSA PUBLIC KEY-----
%s
%s-----END RSA PUBLIC KEY-----`, jwk.N, jwk.E)))
}
