package postgres

import (
	"database/sql"
	"fmt"

	"market_system/models"

	"github.com/google/uuid"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{
		db: db,
	}
}
func (r *orderRepo) Create(req *models.CreateOrder) (*models.Order, error) {
	orderID := uuid.New().String()
	deliveryPrice := 0

	deliveryPriceQuery := `
		SELECT 
		COALESCE(delivery_price, 0)
		FROM
			branches
		WHERE id = $1
	`

	err := r.db.QueryRow(deliveryPriceQuery, req.BranchId).Scan(&deliveryPrice)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO "order"(
			"id",
			"order_id",
			"client_id",
			"branch_id",
			"address",
			"delivery_price",
			"total_count",
			"total_price",
			"updated_at"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  NOW())
	`

	_, err = r.db.Exec(
		query,
		orderID,
		req.OrderId,
		req.ClientId,
		req.BranchId,
		req.Address,
		deliveryPrice,
		req.TotalCount,
		req.TotalPrice,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.OrderPrimaryKey{Id: orderID})
}

func (r *orderRepo) GetByID(req *models.OrderPrimaryKey) (*models.Order, error) {
	var (
		order models.Order
		query = `
			SELECT
				"id",
				"order_id",
				"client_id",
				"branch_id",
				"address",
				"delivery_price",
				COALESCE("total_count",0),
				COALESCE("total_price",0),
				"status",
				"created_at",
				"updated_at"
			FROM "order"
			WHERE "id" = $1
		`
	)

	err := r.db.QueryRow(query, req.Id).Scan(
		&order.Id,
		&order.OrderId,
		&order.ClientId,
		&order.BranchId,
		&order.Address,
		&order.DeliveryPrice,
		&order.TotalCount,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepo) GetList(req *models.GetListOrderRequest) (*models.GetListOrderResponse, error) {
	var (
		resp   models.GetListOrderResponse
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
		where += fmt.Sprintf(`
			AND (
				client_id = '%s'
				OR branch_id = '%s'
			)`, req.Search, req.Search)
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id", 
			"order_id",
			"client_id",
			"branch_id",
			"address",
			"delivery_price",
			COALESCE("total_count", 0),
			COALESCE("total_price", 0),
			"status",
			"created_at",
			"updated_at"
		FROM "order"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order

		err = rows.Scan(
			&resp.Count,
			&order.Id,
			&order.OrderId,
			&order.ClientId,
			&order.BranchId,
			&order.Address,
			&order.DeliveryPrice,
			&order.TotalCount,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Orders = append(resp.Orders, &order)
	}

	return &resp, nil
}

func (r *orderRepo) Update(req *models.UpdateOrder) (int64, error) {
	query := `
		UPDATE "order"
			SET
				client_id = $2,
				branch_id = $3,
				address = $4,
				delivery_price = $5,
				total_count = $6,
				total_price = $7,
				status = $8,
				updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		req.Id,
		req.ClientId,
		req.BranchId,
		req.Address,
		req.DeliveryPrice,
		req.TotalCount,
		req.TotalPrice,
		req.Status,
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

func (r *orderRepo) Delete(req *models.OrderPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM \"order\" WHERE id = $1", req.Id)
	return err
}

// status check

func (r *orderRepo) StatusUpdate(req models.CheckStatus) (models.Order, error) {
	var order models.Order

	selectQuery := `
		SELECT
			"id",
			"order_id",
			"client_id",
			"branch_id",
			"address",
			"delivery_price",
			COALESCE("total_count", 0),
			COALESCE("total_price", 0),
			"status",
			"created_at",
			"updated_at"
		FROM "order"
		WHERE "id" = $1
	`

	err := r.db.QueryRow(selectQuery, req.Id).Scan(
		&order.Id,
		&order.OrderId,
		&order.ClientId,
		&order.BranchId,
		&order.Address,
		&order.DeliveryPrice,
		&order.TotalCount,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return models.Order{}, err
	}

	switch order.Status {
	case "new":
		if req.Status == "canceled" || req.Status == "in-process" {
			updateQuery := `
				UPDATE "order"
				SET
					status = $2,
					updated_at = NOW()
				WHERE id = $1
				RETURNING
					"id",
					"order_id",
					"client_id",
					"branch_id",
					"address",
					"delivery_price",
					COALESCE("total_count", 0),
					COALESCE("total_price", 0),
					"status",
					"created_at",
					"updated_at"
			`
			err := r.db.QueryRow(updateQuery, order.Id, req.Status).Scan(
				&order.Id,
				&order.OrderId,
				&order.ClientId,
				&order.BranchId,
				&order.Address,
				&order.DeliveryPrice,
				&order.TotalCount,
				&order.TotalPrice,
				&order.Status,
				&order.CreatedAt,
				&order.UpdatedAt,
			)
			if err != nil {
				return models.Order{}, err
			}
		} else {
			return models.Order{}, fmt.Errorf("invalid status transition from 'new' to '%s'", req.Status)
		}
	case "in-process":
		if req.Status == "finished"  {
			updateQuery := `
				UPDATE "order"
				SET
					status = $2,
					updated_at = NOW()
				WHERE id = $1
				RETURNING
					"id",
					"order_id",
					"client_id",
					"branch_id",
					"address",
					"delivery_price",
					COALESCE("total_count", 0),
					COALESCE("total_price", 0),
					"status",
					"created_at",
					"updated_at"
			`
			err := r.db.QueryRow(updateQuery, order.Id, req.Status).Scan(
				&order.Id,
				&order.OrderId,
				&order.ClientId,
				&order.BranchId,
				&order.Address,
				&order.DeliveryPrice,
				&order.TotalCount,
				&order.TotalPrice,
				&order.Status,
				&order.CreatedAt,
				&order.UpdatedAt,
			)
			if err != nil {
				return models.Order{}, err
			}
		} else {
			return models.Order{}, fmt.Errorf("invalid status transition from 'in-process' to '%s'", req.Status)
		}
	default:
		return models.Order{}, fmt.Errorf("unsupported status: %s", order.Status)
	}

	return order, nil
}
