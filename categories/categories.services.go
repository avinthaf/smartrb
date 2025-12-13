package categories

import "database/sql"

func getAllCategoriesService(db *sql.DB) ([]Category, error) {
	return getAllCategories(db)
}

func getPrimaryCategoriesService(db *sql.DB) ([]Category, error) {
	return getPrimaryCategories(db)
}

func getCategoriesByIdsService(categoryIds []string, db *sql.DB) ([]Category, error) {
	return getCategoriesByIds(categoryIds, db)
}

func getProductCategoriesByProductIdService(productId string, db *sql.DB) ([]ProductCategory, error) {
	return getProductCategoriesByProductId(productId, db)
}

func getProductCategoriesByProductIdsService(productIds []string, db *sql.DB) ([]ProductCategory, error) {
	return getProductCategoriesByProductIds(productIds, db)
}

func createProductCategoryService(db *sql.DB, productID string, categoryID string) error {
	return createProductCategory(db, productID, categoryID)
}

