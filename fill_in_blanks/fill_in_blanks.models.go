package fill_in_blanks

import "time"

type FillInBlankDeck struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	UserId        *string   `json:"user_id,omitempty"`
	PublishStatus string    `json:"publish_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FillInBlank struct {
	Id          string    `json:"id"`
	DeckId      string    `json:"deck_id"`
	Prompt      string    `json:"prompt"`
	Answers     []string  `json:"answers"`
	Explanation string    `json:"explanation"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FillInBlankDeckSession struct {
	Id        string    `json:"id"`
	DeckId    string    `json:"deck_id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FillInBlankScore struct {
	Id            string    `json:"id"`
	FillInBlankId string    `json:"fill_in_blank_id"`
	UserId        string    `json:"user_id"`
	Score         float64   `json:"score"`
	SessionId     string    `json:"session_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
