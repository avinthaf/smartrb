package flashcards

import (
	"database/sql"
	"fmt"
)

func getFlashcardDecks(db *sql.DB) ([]FlashcardDeck, error) {
	query := `
        SELECT 
            id, 
            title, 
            description, 
            user_id, 
            publish_status,
            created_at, 
            updated_at 
        FROM flashcard_decks
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []FlashcardDeck
	for rows.Next() {
		var deck FlashcardDeck
		var (
			description sql.NullString
			userID      sql.NullString
		)

		err := rows.Scan(
			&deck.Id,
			&deck.Title,
			&description,
			&userID,
			&deck.PublishStatus,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle NULL values
		if description.Valid {
			deck.Description = &description.String
		}

		if userID.Valid {
			deck.UserId = &userID.String
		}

		decks = append(decks, deck)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return decks, nil
}

func getFlashcardsByDeckId(db *sql.DB, deckId string) (GetFlashcardsByDeckIdResult, error) {

	deck, err := getFlashcardDeckById(db, deckId)
	if err != nil {
		return GetFlashcardsByDeckIdResult{}, err
	}

	query := `
        SELECT 
            id, 
            term, 
            definition, 
            deck_id,
            created_at, 
            updated_at 
        FROM flashcards
        WHERE deck_id = $1
    `
	rows, err := db.Query(query, deckId)
	if err != nil {
		return GetFlashcardsByDeckIdResult{}, err
	}
	defer rows.Close()

	var flashcards []Flashcard
	for rows.Next() {
		var flashcard Flashcard
		err := rows.Scan(
			&flashcard.Id,
			&flashcard.Term,
			&flashcard.Definition,
			&flashcard.DeckId,
			&flashcard.CreatedAt,
			&flashcard.UpdatedAt,
		)
		if err != nil {
			return GetFlashcardsByDeckIdResult{}, err
		}
		flashcards = append(flashcards, flashcard)
	}

	if err = rows.Err(); err != nil {
		return GetFlashcardsByDeckIdResult{}, err
	}

	return GetFlashcardsByDeckIdResult{
		Deck: deck,
		Flashcards: flashcards,
	}, nil
}

func getFlashcardDeckById(db *sql.DB, deckId string) (FlashcardDeck, error) {
	query := `
        SELECT 
            id, 
            title, 
            description, 
            user_id, 
            publish_status,
            created_at, 
            updated_at 
        FROM flashcard_decks
        WHERE id = $1
    `
	rows, err := db.Query(query, deckId)
	if err != nil {
		return FlashcardDeck{}, err
	}
	defer rows.Close()

	var deck FlashcardDeck
	for rows.Next() {
		var (
			description sql.NullString
			userID      sql.NullString
		)

		err := rows.Scan(
			&deck.Id,
			&deck.Title,
			&description,
			&userID,
			&deck.PublishStatus,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)

		if err != nil {
			return FlashcardDeck{}, err
		}

		// Handle NULL values
		if description.Valid {
			deck.Description = &description.String
		}

		if userID.Valid {
			deck.UserId = &userID.String
		}
	}

	if err = rows.Err(); err != nil {
		return FlashcardDeck{}, err
	}

	return deck, nil
}

func getFlashcardDeckSessionsByUserId(db *sql.DB, userId string) ([]FlashcardDeckSession, error) {
	query := `
        SELECT 
            id, 
            deck_id, 
            user_id,
            created_at, 
            updated_at 
        FROM flashcard_deck_sessions
        WHERE user_id = $1
    `
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []FlashcardDeckSession
	for rows.Next() {
		var session FlashcardDeckSession
		err := rows.Scan(
			&session.Id,
			&session.DeckId,
			&session.UserId,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func getFlashcardScoresBySessionId(db *sql.DB, sessionId string) ([]FlashcardScore, error) {
	query := `
        SELECT 
            id, 
            card_id, 
            user_id, 
            score, 
            session_id,
            created_at, 
            updated_at 
        FROM flashcard_scores
        WHERE session_id = $1
    `
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []FlashcardScore
	for rows.Next() {
		var score FlashcardScore
		err := rows.Scan(
			&score.Id,
			&score.CardId,
			&score.UserId,
			&score.Score,
			&score.SessionId,
			&score.CreatedAt,
			&score.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return scores, nil
}

func createFlashcardDeckSession(db *sql.DB, id string, deckId string, userId string) (FlashcardDeckSession, error) {
	// First, insert the new session
	query := `
        INSERT INTO flashcard_deck_sessions (id, deck_id, user_id)
        VALUES ($1, $2, $3)
        RETURNING id, deck_id, user_id, created_at, updated_at
    `
	
	var session FlashcardDeckSession
	err := db.QueryRow(query, id, deckId, userId).Scan(
		&session.Id,
		&session.DeckId,
		&session.UserId,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return FlashcardDeckSession{}, fmt.Errorf("failed to create flashcard deck session: %v", err)
	}

	return session, nil
}

func createFlashcardScore(db *sql.DB, userId string, cardId string, score int, sessionId string) error {
	query := `
        INSERT INTO flashcard_scores (card_id, user_id, score, session_id)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, cardId, userId, score, sessionId)
	if err != nil {
		return fmt.Errorf("failed to create flashcard score: %v", err)
	}
	return nil
}