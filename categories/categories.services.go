package categories

import "database/sql"

func getPrimaryCategoriesService(db *sql.DB) ([]Category, error) {
	return getPrimaryCategories(db)
}

func getProductsCategoriesByProductIdService(productId string, db *sql.DB) ([]ProductsCategories, error) {
	return getProductsCategoriesByProductId(productId, db)
}

func getProductsCategoriesByProductIdsService(productIds []string, db *sql.DB) ([]ProductsCategories, error) {
	return getProductsCategoriesByProductIds(productIds, db)
}

