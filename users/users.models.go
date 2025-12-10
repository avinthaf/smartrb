package users

import "time"

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	ExternalId string `json:"external_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

