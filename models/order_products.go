package models

type OrderProductPrimaryKey struct {
	OrderProductID string `json:"order_product_id"`
}

type CreateOrderProduct struct {
	OrderID        string  `json:"order_id"`
	ProductID      string  `json:"product_id"`
	DiscountType   string  `json:"discount_type"`
	DiscountAmount float64 `json:"discount_amount"`
	Quantity       float64 `json:"quantity"`
	Price          float64 `json:"price"`
	Sum            float64 `json:"sum"`
}

type OrderProduct struct {
	OrderProductID string  `json:"order_product_id"`
	OrderID        string  `json:"order_id"`
	ProductID      string  `json:"product_id"`
	DiscountType   string  `json:"discount_type"`
	DiscountAmount float64 `json:"discount_amount"`
	Quantity       float64 `json:"quantity"`
	Price          float64 `json:"price"`
	Sum            float64 `json:"sum"`
}

type UpdateOrderProduct struct {
	OrderProductID string  `json:"order_product_id"`
	OrderID        string  `json:"order_id"`
	ProductID      string  `json:"product_id"`
	DiscountType   string  `json:"discount_type"`
	DiscountAmount float64 `json:"discount_amount"`
	Quantity       float64 `json:"quantity"`
	Price          float64 `json:"price"`
	Sum            float64 `json:"sum"`
}

type GetListOrderProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListOrderProductResponse struct {
	Count         int             `json:"count"`
	OrderProducts []*OrderProduct `json:"order_products"`
}
