package categories

import "database/sql"


type MqCallback func(topic string, routingKey string, message string)

func GetPrimaryCategories(db *sql.DB) ([]Category, error) {
	return getPrimaryCategoriesService(db)
}

func GetProductsCategoriesByProductId(productId string, db *sql.DB) ([]ProductsCategories, error) {
	return getProductsCategoriesByProductIdService(productId, db)
}

func GetProductsCategoriesByProductIds(productIds []string, db *sql.DB) ([]ProductsCategories, error) {
	return getProductsCategoriesByProductIdsService(productIds, db)
}