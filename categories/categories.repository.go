package categories

import (
	"database/sql"
	"fmt"
	"strings"
)

func getAllCategories(db *sql.DB) ([]Category, error) {
	query := `
	SELECT 
		id, 
		name, 
		description, 
		COALESCE(image_url, '') as image_url, 
		parent_id,
		created_at, 
		updated_at 
	FROM categories
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		var parentId *string // Use pointer to handle NULL values

		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.ImageUrl,
			&parentId,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable parent_id to string
		if parentId != nil {
			category.ParentId = *parentId
		} else {
			category.ParentId = ""
		}

		categories = append(categories, category)
	}

	return categories, nil
}

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

func getCategoriesByIds(categoryIds []string, db *sql.DB) ([]Category, error) {
	// Convert string slice to postgres array format
	placeholders := make([]string, len(categoryIds))
	args := make([]interface{}, len(categoryIds))

	for i, id := range categoryIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT 
            id, 
            name, 
            description, 
            COALESCE(image_url, '') as image_url, 
            created_at, 
            updated_at 
        FROM categories 
        WHERE id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
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

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func getProductCategoriesByProductId(productId string, db *sql.DB) ([]ProductCategory, error) {
	query := `
        SELECT 
            id, 
            product_id, 
            category_id, 
            created_at, 
            updated_at 
        FROM product_categories 
        WHERE product_id = $1`

	rows, err := db.Query(query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsCategories []ProductCategory
	for rows.Next() {
		var productCategory ProductCategory

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

func getProductCategoriesByProductIds(productIds []string, db *sql.DB) ([]ProductCategory, error) {
	// Convert string slice to postgres array format
	placeholders := make([]string, len(productIds))
	args := make([]interface{}, len(productIds))

	for i, id := range productIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT 
            id, 
            product_id, 
            category_id, 
            created_at
        FROM product_categories 
        WHERE product_id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsCategories []ProductCategory
	for rows.Next() {
		var productCategory ProductCategory

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productsCategories, nil
}

func createProductCategory(db *sql.DB, productID string, categoryID string) error {
	query := `
        INSERT INTO product_categories (product_id, category_id, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())`

	_, err := db.Exec(query, productID, categoryID)
	if err != nil {
		return err
	}
	return nil
}
