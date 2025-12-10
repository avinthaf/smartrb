package flashcards

import "database/sql"

func getFlashcardDecksService(db *sql.DB) ([]FlashcardDeck, error) {
	return getFlashcardDecks(db)
}

func getFlashcardsByDeckIdService(db *sql.DB, deckId string) (GetFlashcardsByDeckIdResult, error) {
	return getFlashcardsByDeckId(db, deckId)
}

func createFlashcardDeckSessionService(db *sql.DB, id string, deckId string, userId string) (FlashcardDeckSession, error) {
	return createFlashcardDeckSession(db, id, deckId, userId)
}

func createFlashcardScoreService(db *sql.DB, userId string, cardId string, score int, sessionId string) error {
	return createFlashcardScore(db, userId, cardId, score, sessionId)
}

