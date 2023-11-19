package models

type ClientPrimaryKey struct {
	ID string `json:"id"`
}

type CreateClient struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	DateOfBirth string `json:"date_of_birth"`
	CreatedAt  string `json:"created_at"`
}

type Client struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	DateOfBirth string `json:"date_of_birth"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UpdateClient struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	DateOfBirth string `json:"date_of_birth"`
}

type GetListClientRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListClientResponse struct {
	Count   int      `json:"count"`
	Clients []*Client `json:"clients"`
}
