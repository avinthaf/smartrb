package products

import "database/sql"

func createProductRating(db *sql.DB, productId string, userId string, rating int8) error {
	
	query := `INSERT INTO product_ratings (product_id, user_id, rating) VALUES ($1, $2, $3)`
	
	_, err := db.Exec(query, productId, userId, rating)
	if err != nil {
		return err
	}
	return nil
}