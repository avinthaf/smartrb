package products

import "database/sql"

type CreateProductRatingRequest struct {
	ProductId string `json:"productId"`
	Rating int8 `json:"rating"`
}

func CreateProductRating(db *sql.DB, productId string, userId string, rating int8) error {
	return createProductRatingService(db, productId, userId, rating)
}