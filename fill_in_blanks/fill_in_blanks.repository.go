package fill_in_blanks

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

func getFillInBlankDecks(db *sql.DB) ([]FillInBlankDeck, error) {
	query := `
        SELECT 
            id, 
            title, 
            description, 
            user_id, 
            publish_status,
            created_at, 
            updated_at 
        FROM fill_in_blank_decks
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []FillInBlankDeck
	for rows.Next() {
		var deck FillInBlankDeck
		var userID sql.NullString

		err := rows.Scan(
			&deck.Id,
			&deck.Title,
			&deck.Description,
			&userID,
			&deck.PublishStatus,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle NULL values
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

func getFillInBlanksByDeckId(db *sql.DB, deckId string) (GetFillInBlanksByDeckIdResult, error) {

	deck, err := getFillInBlankDeckById(db, deckId)
	if err != nil {
		return GetFillInBlanksByDeckIdResult{}, err
	}

	query := `
        SELECT 
            id, 
            deck_id,
            prompt, 
            answers,
            explanation,
            created_at, 
            updated_at 
        FROM fill_in_blanks
        WHERE deck_id = $1
    `
	rows, err := db.Query(query, deckId)
	if err != nil {
		return GetFillInBlanksByDeckIdResult{}, err
	}
	defer rows.Close()

	var fillInBlanks []FillInBlank
	for rows.Next() {
		var fillInBlank FillInBlank
		var answers pq.StringArray
		
		err := rows.Scan(
			&fillInBlank.Id,
			&fillInBlank.DeckId,
			&fillInBlank.Prompt,
			&answers,
			&fillInBlank.Explanation,
			&fillInBlank.CreatedAt,
			&fillInBlank.UpdatedAt,
		)
		if err != nil {
			return GetFillInBlanksByDeckIdResult{}, err
		}

		// Convert pq.StringArray to []string
		fillInBlank.Answers = []string(answers)

		fillInBlanks = append(fillInBlanks, fillInBlank)
	}

	if err = rows.Err(); err != nil {
		return GetFillInBlanksByDeckIdResult{}, err
	}

	return GetFillInBlanksByDeckIdResult{
		Deck:         deck,
		FillInBlanks: fillInBlanks,
	}, nil
}

func getFillInBlankDeckById(db *sql.DB, deckId string) (FillInBlankDeck, error) {
	query := `
        SELECT 
            id, 
            title, 
            description, 
            user_id, 
            publish_status,
            created_at, 
            updated_at 
        FROM fill_in_blank_decks
        WHERE id = $1
    `
	rows, err := db.Query(query, deckId)
	if err != nil {
		return FillInBlankDeck{}, err
	}
	defer rows.Close()

	var deck FillInBlankDeck
	for rows.Next() {
		var userID sql.NullString

		err := rows.Scan(
			&deck.Id,
			&deck.Title,
			&deck.Description,
			&userID,
			&deck.PublishStatus,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)

		if err != nil {
			return FillInBlankDeck{}, err
		}

		// Handle NULL values
		if userID.Valid {
			deck.UserId = &userID.String
		}
	}

	if err = rows.Err(); err != nil {
		return FillInBlankDeck{}, err
	}

	return deck, nil
}

func getFillInBlankDeckSessionsByUserId(db *sql.DB, userId string) ([]FillInBlankDeckSession, error) {
	query := `
        SELECT 
            id, 
            deck_id, 
            user_id,
            created_at, 
            updated_at 
        FROM fill_in_blank_deck_sessions
        WHERE user_id = $1
    `
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []FillInBlankDeckSession
	for rows.Next() {
		var session FillInBlankDeckSession
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

func getFillInBlankScoresBySessionId(db *sql.DB, sessionId string) ([]FillInBlankScore, error) {
	query := `
        SELECT 
            id, 
            fill_in_blank_id, 
            user_id, 
            score, 
            session_id,
            created_at, 
            updated_at 
        FROM fill_in_blank_scores
        WHERE session_id = $1
    `
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []FillInBlankScore
	for rows.Next() {
		var score FillInBlankScore
		err := rows.Scan(
			&score.Id,
			&score.FillInBlankId,
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

func createFillInBlankDeckSession(db *sql.DB, id string, deckId string, userId string) (FillInBlankDeckSession, error) {
	// First, insert the new session
	query := `
        INSERT INTO fill_in_blank_deck_sessions (id, deck_id, user_id)
        VALUES ($1, $2, $3)
        RETURNING id, deck_id, user_id, created_at, updated_at
    `
	
	var session FillInBlankDeckSession
	err := db.QueryRow(query, id, deckId, userId).Scan(
		&session.Id,
		&session.DeckId,
		&session.UserId,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return FillInBlankDeckSession{}, fmt.Errorf("failed to create fill in blank deck session: %v", err)
	}

	return session, nil
}

func createFillInBlankScore(db *sql.DB, userId string, fillInBlankId string, score float64, sessionId string) error {
	query := `
        INSERT INTO fill_in_blank_scores (fill_in_blank_id, user_id, score, session_id)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, fillInBlankId, userId, score, sessionId)
	if err != nil {
		return fmt.Errorf("failed to create fill in blank score: %v", err)
	}
	return nil
}