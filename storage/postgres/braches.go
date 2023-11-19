package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"market_system/models"

	"github.com/google/uuid"
)

type branchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(req *models.CreateBranch) (*models.Branch, error) {
	branchID := uuid.New().String()
	query := `
		INSERT INTO "branches"(
			"id",
			"name", 
			"phone", 
			"photo",
			"work_start_hour",
			"work_end_hour",
			"address",
			"created_at",
			"updated_at"
		) VALUES ($1, $2, $3, $4, $5, $6, $7,  NOW(), NOW())
	`

	_, err := r.db.Exec(
		query,
		branchID,
		req.Name,
		req.Phone,
		req.Photo,
		req.WorkStartHour,
		req.WorkEndHour,
		req.Address,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.BranchPrimaryKey{ID: branchID})
}

func (r *branchRepo) GetByID(req *models.BranchPrimaryKey) (*models.Branch, error) {
	var (
		branch models.Branch
		query  = `
			SELECT
				"id",
				"name",
				"phone",
				"photo",
				"work_start_hour",
				"work_end_hour",
				"address",
				"delivery_price",
				"active",
				"created_at",
				"updated_at"	
			FROM "branches"
			WHERE "id" = $1
		`
	)

	err := r.db.QueryRow(query, req.ID).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Phone,
		&branch.Photo,
		&branch.WorkStartHour,
		&branch.WorkEndHour,
		&branch.Address,
		&branch.DeliveryPrice,
		&branch.Active,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &branch, nil
}

func (r *branchRepo) Update(req *models.UpdateBranch) (int64, error) {
	query := `
		UPDATE branches
			SET
				"name" = $2,
				"phone" = $3,
				"photo" = $4,
				"work_start_hour" = $5,
				"work_end_hour" = $6,
				"address" = $7,
				"active" = $8,
				"updated_at" = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		req.ID,
		req.Name,
		req.Phone,
		req.Photo,
		req.WorkStartHour,
		req.WorkEndHour,
		req.Address,
		req.Active,
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

func (r *branchRepo) Delete(req *models.BranchPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM branches WHERE id = $1", req.ID)
	return err
}

func (r *branchRepo) GetList(req *models.GetListBranchRequest) (*models.GetListBranchResponse, error) {
	var (
		resp   models.GetListBranchResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	hour := strconv.Itoa(time.Now().Hour())
	minute := strconv.Itoa(time.Now().Minute())

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	currentTime := hour + ":" + minute
	if len(req.Search) > 0 {
		where += " AND name ILIKE '%" + req.Search + "%' AND '" + currentTime + "' BETWEEN work_start_hour AND work_end_hour"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"phone",
			"photo",
			"work_start_hour",
			"work_end_hour",
			"address",
			"delivery_price",
			"active",
			"created_at",
			"updated_at"
		FROM "branches"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branch models.Branch

		err = rows.Scan(
			&resp.Count,
			&branch.ID,
			&branch.Name,
			&branch.Phone,
			&branch.Photo,
			&branch.WorkStartHour,
			&branch.WorkEndHour,
			&branch.Address,
			&branch.DeliveryPrice,
			&branch.Active,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if !isBranchOpen(&branch, currentTime) {
			branch.Active = false
		}

		resp.Branches = append(resp.Branches, &branch)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}

// isBranchOpen проверяет, открыт ли бранч в текущее время
func isBranchOpen(branch *models.Branch, currentTime string) bool {
	return currentTime >= branch.WorkStartHour && currentTime <= branch.WorkEndHour
}
