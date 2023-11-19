package postgres

import (
	"database/sql"
	"fmt"

	"market_system/models"

	"github.com/google/uuid"
)

type clientRepo struct {
	db *sql.DB
}

func NewClientRepo(db *sql.DB) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (r *clientRepo) Create(req *models.CreateClient) (*models.Client, error) {
	clientID := uuid.New().String()
	query := `
		INSERT INTO "client"(
			"id",
			"first_name", 
			"last_name", 
			"phone",
			"photo",
			"date_of_birth",
			"created_at",
			"updated_at"
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.Exec(
		query,
		clientID,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Photo,
		req.DateOfBirth,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.ClientPrimaryKey{ID: clientID})
}

func (r *clientRepo) GetByID(req *models.ClientPrimaryKey) (*models.Client, error) {
	var (
		client models.Client
		query  = `
			SELECT
				"id",
				"first_name",
				"last_name",
				"phone",
				"photo",
				"date_of_birth",
				"created_at",
				"updated_at"	
			FROM "client"
			WHERE "id" = $1
		`
	)

	err := r.db.QueryRow(query, req.ID).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.Phone,
		&client.Photo,
		&client.DateOfBirth,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}
func (r *clientRepo) GetList(req *models.GetListClientRequest) (*models.GetListClientResponse, error) {
	var (
		resp   models.GetListClientResponse
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
		where += " AND first_name ILIKE" + " '%" + req.Search + "%'" + "OR last_name ILIKE" + "'%" + req.Search + "%'" + "OR phone ILIKE" + "'%" + req.Search + "%'"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"first_name",
			"last_name",
			"phone",
			"photo",
			"date_of_birth",
			"created_at",
			"updated_at"
		FROM "client"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var client models.Client

		err = rows.Scan(
			&resp.Count,
			&client.ID,
			&client.FirstName,
			&client.LastName,
			&client.Phone,
			&client.Photo,
			&client.DateOfBirth,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Clients = append(resp.Clients, &client)
	}

	return &resp, nil
}

func (r *clientRepo) Update(req *models.UpdateClient) (int64, error) {
	query := `
		UPDATE client
			SET
				first_name = $2,
				last_name = $3,
				phone = $4,
				photo = $5,
				date_of_birth = $6,
				updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		req.ID,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Photo,
		req.DateOfBirth,
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

func (r *clientRepo) Delete(req *models.ClientPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM client WHERE id = $1", req.ID)
	return err
}
