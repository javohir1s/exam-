package models

type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateProduct struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductId   string `json:"product_id"`
	CategoryId  string `json:"category_id"`
	Photo       string `json:"photo"`
	Price       float64 `json:"price"`
}

type Product struct {
	Id          string      `json:"id"`
	ProductId   string      `json:"product_id"`
	Title       string      `json:"title"`
	Photo       string      `json:"photo"`
	Price       float64      `json:"price"`
	Description string      `json:"description"`
	CategoryId  string      `json:"category_id"`
	Category    interface{} `json:"category"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type UpdateProduct struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductId   string `json:"product_id"`
	Photo       string `json:"photo"`
	Price       float64 `json:"price"`
	CategoryId  string `json:"category_id"`
}

type GetListProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
}

type GetListProductResponse struct {
	Count    int        `json:"count"`
	Products []*Product `json:"products"`
}
