package models

type Order struct {
	Id            string  `json:"id"`
	OrderId       string  `json:"order_id"`
	ClientId      string  `json:"client_id"`
	BranchId      string  `json:"branch_id"`
	Address       string  `json:"address"`
	DeliveryPrice float64 `json:"delivery_price"`
	TotalCount    float64 `json:"total_count"`
	TotalPrice    float64 `json:"total_price"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type CreateOrder struct {
	OrderId    string  `json:"order_id"`
	ClientId   string  `json:"client_id"`
	BranchId   string  `json:"branch_id"`
	Address    string  `json:"address"`
	TotalCount float64 `json:"total_count"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type UpdateOrder struct {
	Id            string  `json:"id"`
	ClientId      string  `json:"client_id"`
	OrderId       string  `json:"order_id"`
	BranchId      string  `json:"branch_id"`
	Address       string  `json:"address"`
	DeliveryPrice float64 `json:"delivery_price"`
	TotalCount    float64 `json:"total_count"`
	TotalPrice    float64 `json:"total_price"`
	Status        string  `json:"status"`
}

type GetListOrderRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
}

type GetListOrderResponse struct {
	Count  int      `json:"count"`
	Orders []*Order `json:"orders"`
}

type CheckStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
