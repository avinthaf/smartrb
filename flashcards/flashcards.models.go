package flashcards

import "time"

type FlashcardDeck struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description,omitempty"`
	UserId        *string   `json:"user_id,omitempty"`
	PublishStatus string    `json:"publish_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Flashcard struct {
	Id         string    `json:"id"`
	Term       string    `json:"term"`
	Definition string    `json:"definition"`
	DeckId     string    `json:"deck_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FlashcardDeckSession struct {
	Id        string    `json:"id"`
	DeckId    string    `json:"deck_id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FlashcardScore struct {
	Id        string    `json:"id"`
	CardId    string    `json:"card_id"`
	UserId    string    `json:"user_id"`
	Score     int       `json:"score"` // can only be 0 or 1
	SessionId string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
