package postgres

import (
	"database/sql"
	"fmt"

	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/google/uuid"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(req *models.CreateProduct) (*models.Product, error) {

	query := `
            INSERT INTO "product"(
                "id",
				"product_id",
                "title",
                "description",
                "price",
                "photo",
                "category_id",
                "created_at",
				"updated_at"
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`
	productId := uuid.New().String()

	_, err := r.db.Exec(
		query,
		productId,
		req.ProductId,
		req.Title,
		req.Description,
		req.Price,
		req.Photo,
		helpers.NewNullString(req.CategoryId),
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании продукта: %v", err)
	}

	return r.GetByID(&models.ProductPrimaryKey{Id: productId})
}

func (r *productRepo) GetByID(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query = `
    SELECT
        "id",
        "title",   
		"product_id",
		"description",
        "photo",   
        "price",  
        "category_id",
        "updated_at",
        "created_at"
    FROM "product"
    WHERE "id" = $1
    `
	)

	var (
		id          sql.NullString
		product_id  sql.NullString
		description sql.NullString
		title       sql.NullString
		photo       sql.NullString
		price       sql.NullFloat64
		category_id sql.NullString
		updated_at  sql.NullString
		created_at  sql.NullString
	)

	err := r.db.QueryRow(query, req.Id).Scan(
		&id,
		&title,
		&product_id,
		&description,
		&photo,
		&price,
		&category_id,
		&updated_at,
		&created_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         id.String,
		Title:      title.String,
		Photo:      photo.String,
		ProductId:  product_id.String,
		Description: description.String,
		Price:      price.Float64,
		CategoryId: category_id.String,
		UpdatedAt:  updated_at.String,
		CreatedAt:  created_at.String,
	}, nil
}

func (r *productRepo) GetList(req *models.GetListProductRequest) (*models.GetListProductResponse, error) {
	var (
		resp   models.GetListProductResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if len(req.Search) > 0 {
		where += " AND (title ILIKE '%" + req.Search + "%' OR category_id::text ILIKE '%" + req.Search + "%')"
	}
	
	
	var countQuery = `
		SELECT COUNT(*) FROM "product"
	`
	countQuery += where

	err := r.db.QueryRow(countQuery).Scan(&resp.Count)
	if err != nil {
		return nil, err
	}

	var selectQuery = `
		SELECT
			"id",
			"title",   
			"product_id",
			"description",
			"photo",   
			"price",
			"category_id",
			"updated_at",
			"created_at"
		FROM "product"
	`

	selectQuery += where + sort + offset + limit
	rows, err := r.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			title       sql.NullString
			photo       sql.NullString
			productId   sql.NullString
			description sql.NullString
			price       sql.NullFloat64
			category_id sql.NullString
			updated_at  sql.NullString
			created_at  sql.NullString
		)

		err := rows.Scan(
			&id,
			&title,
			&photo,
			&productId,
			&description,
			&price,
			&category_id,
			&updated_at,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, &models.Product{
			Id:         id.String,
			Title:      title.String,
			Photo:      photo.String,
			Price:      price.Float64,
			CategoryId: category_id.String,
			UpdatedAt:  updated_at.String,
			CreatedAt:  created_at.String,
		})
	}

	return &resp, nil
}

func (r *productRepo) Update(req *models.UpdateProduct) (int64, error) {
	query := `
		UPDATE product
			SET
				title = $2,    
				photo = $3,  
				description = $4,
				price = $5,
				category_id = $6,
				updated_at = NOW()  
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.Id,
		req.Title,
		req.Photo,
		req.Description,
		req.Price,
		helpers.NewNullString(req.CategoryId),
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *productRepo) Delete(req *models.ProductPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM product WHERE id = $1", req.Id)
	return err
}
