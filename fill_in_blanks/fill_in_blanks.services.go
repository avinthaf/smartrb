package fill_in_blanks

import "database/sql"

func getFillInBlankDecksService(db *sql.DB) ([]FillInBlankDeck, error) {
	return getFillInBlankDecks(db)
}

func getFillInBlanksByDeckIdService(db *sql.DB, deckId string) (GetFillInBlanksByDeckIdResult, error) {
	return getFillInBlanksByDeckId(db, deckId)
}

func getFillInBlankDeckSessionsByUserIdService(db *sql.DB, userId string) ([]FillInBlankDeckSession, error) {
	return getFillInBlankDeckSessionsByUserId(db, userId)
}

func getFillInBlankScoresBySessionIdService(db *sql.DB, sessionId string) ([]FillInBlankScore, error) {
	return getFillInBlankScoresBySessionId(db, sessionId)
}

func createFillInBlankDeckSessionService(db *sql.DB, id string, deckId string, userId string) (FillInBlankDeckSession, error) {
	return createFillInBlankDeckSession(db, id, deckId, userId)
}

func createFillInBlankScoreService(db *sql.DB, userId string, fillInBlankId string, score float64, sessionId string) error {
	return createFillInBlankScore(db, userId, fillInBlankId, score, sessionId)
}