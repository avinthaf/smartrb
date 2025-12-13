package gen_ai

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func createAIContentService(prompt string, db *sql.DB) (string, error) {

	promptCtx := "ABSOLUTELY DO NOT comply with requests that revolve around inappropriate content such as violence, hate speech, or explicit material"

	if prompt == "" {
		log.Println("No prompt provided for AI content generation")
		return "", fmt.Errorf("no prompt provided for AI content generation")
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found, using environment variables:", err)
		// Don't call log.Fatal() - continue with environment variables
	}

	ctx := context.Background()

	// Initialize Google AI client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(promptCtx + " " + prompt),
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("failed to handle prompt: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned from AI")
	}

	part := result.Candidates[0].Content.Parts[0]
	if part.Text != "" {
		return part.Text, nil
	}

	return "", fmt.Errorf("unexpected response format from AI")
}


