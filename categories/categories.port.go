package categories

import "database/sql"


type MqCallback func(topic string, routingKey string, message string)

func GetAllCategories(db *sql.DB) ([]Category, error){
	return getAllCategoriesService(db)
}

func GetPrimaryCategories(db *sql.DB) ([]Category, error) {
	return getPrimaryCategoriesService(db)
}

func GetCategoriesByIds(categoryIds []string, db *sql.DB) ([]Category, error) {
	return getCategoriesByIdsService(categoryIds, db)
}

func GetProductCategoriesByProductId(productId string, db *sql.DB) ([]ProductCategory, error) {
	return getProductCategoriesByProductIdService(productId, db)
}

func GetProductCategoriesByProductIds(productIds []string, db *sql.DB) ([]ProductCategory, error) {
	return getProductCategoriesByProductIdsService(productIds, db)
}

func CreateProductCategory(db *sql.DB, productID string, categoryID string) error {
	return createProductCategoryService(db, productID, categoryID)
}