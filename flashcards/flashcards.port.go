package flashcards

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type GetFlashcardsByDeckIdResult struct {
	Deck       FlashcardDeck `json:"deck"`
	Flashcards []Flashcard   `json:"flashcards"`
}

type FlashcardScoreRequest struct {
	CardId    string `json:"card_id"`
	Score     int    `json:"score"`
	SessionId string `json:"session_id"`
}

type FlashcardDeckSessionRequest struct {
	DeckId    string `json:"deck_id"`
	SessionId string `json:"session_id"`
}

type CreateFlashcardDeckRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	UserId      string      `json:"user_id"`
	PublishStatus string    `json:"publish_status"`
	Flashcards  []Flashcard `json:"flashcards,omitempty"`
}

type CreateFlashcardsRequest struct {
	DeckId    string      `json:"deck_id"`
	Flashcards []Flashcard `json:"flashcards"`
}

func GetFlashcardDecks(db *sql.DB) ([]FlashcardDeck, error) {
	return getFlashcardDecksService(db)
}

func GetFlashcardsByDeckId(db *sql.DB, deckId string) (GetFlashcardsByDeckIdResult, error) {
	return getFlashcardsByDeckIdService(db, deckId)
}

func GetFlashcardDeckSessionsByUserId(db *sql.DB, userId string) ([]FlashcardDeckSession, error) {
	return getFlashcardDeckSessionsByUserIdService(db, userId)
}

func GetFlashcardScoresBySessionId(db *sql.DB, sessionId string) ([]FlashcardScore, error) {
	return getFlashcardScoresBySessionIdService(db, sessionId)
}

func CreateFlashcardDeck(db *sql.DB, title string, description string, userId string, publishStatus string) (FlashcardDeck, error) {
	return createFlashcardDeckService(db, title, description, userId, publishStatus)
}

func CreateFlashcards(db *sql.DB, request CreateFlashcardsRequest) error {
	return createFlashcardsService(db, request.DeckId, request.Flashcards)
}

func CreateFlashcardDeckSession(db *sql.DB, id string, deckId string, userId string) (FlashcardDeckSession, error) {
	return createFlashcardDeckSessionService(db, id, deckId, userId)
}

func CreateFlashcardScore(db *sql.DB, userId string, cardId string, score int, sessionId string) error {
	// Validate UUIDs
	if _, err := uuid.Parse(userId); err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}
	if _, err := uuid.Parse(cardId); err != nil {
		return fmt.Errorf("invalid card ID format: %v", err)
	}

	// Validate session ID is a valid UUID if provided
	if sessionId != "" {
		if _, err := uuid.Parse(sessionId); err != nil {
			return fmt.Errorf("invalid session ID format: %v", err)
		}
	}

	// Validate score is between 0 and 1
	if score < 0 || score > 1 {
		return fmt.Errorf("score must be 0 or 1")
	}

	return createFlashcardScoreService(db, userId, cardId, score, sessionId)
}