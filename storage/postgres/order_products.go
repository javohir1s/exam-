package postgres

import (
	"database/sql"
	"fmt"

	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/google/uuid"
)

type orderProductRepo struct {
	db *sql.DB
}

func NewOrderProductRepo(db *sql.DB) *orderProductRepo {
	return &orderProductRepo{
		db: db,
	}
}

func (r *orderProductRepo) Create(req *models.CreateOrderProduct) (*models.OrderProduct, error) {
	orderProductID := uuid.New().String()

	query := `
        INSERT INTO "order_products"(
            "order_product_id",
            "product_id",
			"order_id",
            "discount_type",
            "discount_amount",
            "quantity",
            "price",
            "sum",
            "created_at",
            "updated_at"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7,$8, DEFAULT, DEFAULT)
    `

	_, err := r.db.Exec(
		query,
		orderProductID,
		req.ProductID,

		helpers.NewNullString(req.OrderID),
		req.DiscountType,
		req.DiscountAmount,
		req.Quantity,
		req.Price,
		req.Sum,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.OrderProductPrimaryKey{OrderProductID: orderProductID})
}

func (r *orderProductRepo) GetByID(req *models.OrderProductPrimaryKey) (*models.OrderProduct, error) {
	var (
		orderProduct models.OrderProduct
		query        = `
            SELECT
                "order_product_id",
                "product_id",
                "discount_type",
                "discount_amount",
                "quantity",
                "price",
                "sum"
            FROM "order_products"
            WHERE "order_product_id" = $1
        `
	)

	err := r.db.QueryRow(query, req.OrderProductID).Scan(
		&orderProduct.OrderProductID,
		&orderProduct.ProductID,
		&orderProduct.DiscountType,
		&orderProduct.DiscountAmount,
		&orderProduct.Quantity,
		&orderProduct.Price,
		&orderProduct.Sum,
	)

	if err != nil {
		return nil, err
	}

	return &orderProduct, nil
}

func (r *orderProductRepo) GetList(req *models.GetListOrderProductRequest) (*models.GetListOrderProductResponse, error) {
	var (
		resp   models.GetListOrderProductResponse
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

	var query = `
        SELECT
            COUNT(*) OVER(),
            "order_product_id",
            "product_id",
            "discount_type",
            "discount_amount",
            "quantity",
            "price",
            "sum"
        FROM "order_products"
    `

	query += sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var orderProduct models.OrderProduct

		err = rows.Scan(
			&resp.Count,
			&orderProduct.OrderProductID,
			&orderProduct.ProductID,
			&orderProduct.DiscountType,
			&orderProduct.DiscountAmount,
			&orderProduct.Quantity,
			&orderProduct.Price,
			&orderProduct.Sum,
		)
		if err != nil {
			return nil, err
		}

		resp.OrderProducts = append(resp.OrderProducts, &orderProduct)
	}

	return &resp, nil
}

func (r *orderProductRepo) Update(req *models.UpdateOrderProduct) (int64, error) {
	query := `
        UPDATE "order_products"
            SET
                "product_id" = $2,
                "discount_type" = $3,
                "discount_amount" = $4,
                "quantity" = $5,
                "price" = $6,
                "sum" = $7,
                "updated_at" = DEFAULT
        WHERE "order_product_id" = $1
    `

	result, err := r.db.Exec(
		query,
		req.OrderProductID,
		req.ProductID,
		req.DiscountType,
		req.DiscountAmount,
		req.Quantity,
		req.Price,
		req.Sum,
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

func (r *orderProductRepo) Delete(req *models.OrderProductPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM \"order_products\" WHERE \"order_product_id\" = $1", req.OrderProductID)
	return err
}


