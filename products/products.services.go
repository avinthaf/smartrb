package products

import "database/sql"

func createProductRatingService(db *sql.DB, productId string, userId string, rating int8) error {
	return createProductRating(db, productId, userId, rating)
}