package models

type BranchPrimaryKey struct {
	ID string `json:"id"`
}

type CreateBranch struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Photo         string `json:"photo"`
	WorkStartHour string `json:"work_start_hour"`
	WorkEndHour   string `json:"work_end_hour"`
	Address       string `json:"address"`
	DeliveryPrice float64  `json:"delivery_price"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
}

type Branch struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Photo         string `json:"photo"`
	WorkStartHour string `json:"work_start_hour"`
	WorkEndHour   string `json:"work_end_hour"`
	Address       string `json:"address"`
	DeliveryPrice float64  `json:"delivery_price"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UpdateBranch struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Photo         string `json:"photo"`
	WorkStartHour string `json:"work_start_hour"`
	WorkEndHour   string `json:"work_end_hour"`
	Address       string `json:"address"`
	DeliveryPrice float64  `json:"delivery_price"`
	Active        bool   `json:"active"`
}

type GetListBranchRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Query  string `json:"query"`
}

type GetListBranchResponse struct {
	Count    int       `json:"count"`
	Branches []*Branch `json:"branches"`
}
