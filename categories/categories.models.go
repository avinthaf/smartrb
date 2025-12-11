package categories

import "time"

type Category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ParentId string `json:"parent_id"`
	ImageUrl string `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductCategory struct {
	Id string `json:"id"`
	ProductId string `json:"product_id"`
	CategoryId string `json:"category_id"`
	CreatedAt time.Time `json:"created_at"`
}