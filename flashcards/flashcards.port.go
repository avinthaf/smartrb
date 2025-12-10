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

func GetFlashcardDecks(db *sql.DB) ([]FlashcardDeck, error) {
	return getFlashcardDecksService(db)
}

func GetFlashcardsByDeckId(db *sql.DB, deckId string) (GetFlashcardsByDeckIdResult, error) {
	return getFlashcardsByDeckIdService(db, deckId)
}

func GetFlashcardScoresBySessionId(db *sql.DB, sessionId string) ([]FlashcardScore, error) {
	return getFlashcardScoresBySessionIdService(db, sessionId)
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