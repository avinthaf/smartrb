package interests

import "time"

type Interest struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	CategoryId string    `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
