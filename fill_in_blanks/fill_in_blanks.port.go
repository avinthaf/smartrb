package fill_in_blanks

import (
	"database/sql"
	"crypto/rand"
	"encoding/hex"
)

type GetFillInBlanksByDeckIdResult struct {
	Deck          FillInBlankDeck `json:"deck"`
	FillInBlanks  []FillInBlank   `json:"fill_in_blanks"`
}

type FillInBlankScoreRequest struct {
	FillInBlankId string `json:"fill_in_blank_id"`
	Score         float64 `json:"score"`
	SessionId     string `json:"session_id"`
}

type FillInBlankDeckSessionRequest struct {
	DeckId    string `json:"deck_id"`
	SessionId string `json:"session_id"`
}

func generateId() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GetFillInBlankDecks(db *sql.DB) ([]FillInBlankDeck, error) {
	return getFillInBlankDecksService(db)
}

func GetFillInBlanksByDeckId(db *sql.DB, deckId string) (GetFillInBlanksByDeckIdResult, error) {
	return getFillInBlanksByDeckIdService(db, deckId)
}

func GetFillInBlankDeckSessionsByUserId(db *sql.DB, userId string) ([]FillInBlankDeckSession, error) {
	return getFillInBlankDeckSessionsByUserIdService(db, userId)
}

func GetFillInBlankScoresBySessionId(db *sql.DB, sessionId string) ([]FillInBlankScore, error) {
	return getFillInBlankScoresBySessionIdService(db, sessionId)
}

func CreateFillInBlankDeckSession(db *sql.DB, request FillInBlankDeckSessionRequest, userId string) (FillInBlankDeckSession, error) {
	id := generateId()
	return createFillInBlankDeckSessionService(db, id, request.DeckId, userId)
}

func CreateFillInBlankScore(db *sql.DB, request FillInBlankScoreRequest, userId string) error {
	return createFillInBlankScoreService(db, userId, request.FillInBlankId, request.Score, request.SessionId)
}

