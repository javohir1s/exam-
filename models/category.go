package models

type CategoryPrimaryKey struct {
	Id string `json:"id"`
}

type CreateCategory struct {
	Title     string `json:"title"`
	ParentID  string `json:"parent_id"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

type Category struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	ParentID  string `json:"parent_id"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateCategory struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ParentID string `json:"parent_id"`
	Image    string `json:"image"`
}

type GetListCategoryRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListCategoryResponse struct {
	Count      int         `json:"count"`
	Categories []*Category `json:"categories"`
}
