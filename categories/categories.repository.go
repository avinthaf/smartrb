package categories

import (
	"database/sql"
)

func getPrimaryCategories(db *sql.DB) ([]Category, error) {
	query := `
        SELECT 
            id, 
            name, 
            description, 
            COALESCE(image_url, '') as image_url, 
            created_at, 
            updated_at 
        FROM categories 
        WHERE parent_id IS NULL`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.ImageUrl,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set empty parent_id since we're not selecting it
		category.ParentId = ""

		categories = append(categories, category)
	}

	return categories, nil
}

func getProductsCategoriesByProductId(productId string, db *sql.DB) ([]ProductsCategories, error) {
	query := `
        SELECT 
            id, 
            product_id, 
            category_id, 
            created_at, 
            updated_at 
        FROM products_categories 
        WHERE product_id = $1`

	rows, err := db.Query(query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsCategories []ProductsCategories
	for rows.Next() {
		var productCategory ProductsCategories

		err := rows.Scan(
			&productCategory.Id,
			&productCategory.ProductId,
			&productCategory.CategoryId,
			&productCategory.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		productsCategories = append(productsCategories, productCategory)
	}

	return productsCategories, nil
}

func getProductsCategoriesByProductIds(productIds []string, db *sql.DB) ([]ProductsCategories, error) {
	query := `
        SELECT 
            id, 
            product_id, 
            category_id, 
            created_at, 
            updated_at 
        FROM products_categories 
        WHERE product_id IN ($1)`

	rows, err := db.Query(query, productIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsCategories []ProductsCategories
	for rows.Next() {
		var productCategory ProductsCategories

		err := rows.Scan(
			&productCategory.Id,
			&productCategory.ProductId,
			&productCategory.CategoryId,
			&productCategory.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		productsCategories = append(productsCategories, productCategory)
	}

	return productsCategories, nil
}
