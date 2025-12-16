package products

type ProductRating struct {
	Id string `json:"id"`
	ProductId string `json:"product_id"`
	UserId string `json:"user_id"`
	Rating int8 `json:"rating"` // Can only be between 0 and 5
	CreatedAt string `json:"created_at"`
}