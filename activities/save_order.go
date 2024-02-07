package activities

import (
	"bitovi/distributed-systems-examples/schemas"
	"context"
	"database/sql"
	"encoding/json"
)

type SaveOrderConfig struct {
	DB *sql.DB
}

func NewSaveOrder(db *sql.DB) *SaveOrderConfig {
	return &SaveOrderConfig{
		DB: db,
	}
}

func (cfg SaveOrderConfig) SaveOrderActivity(ctx context.Context, i schemas.WorkflowInput) (*string, error) {
	customerJSON, err := json.Marshal(i.Customer)
	if err != nil {
		return nil, err
	}
	productsJSON, err := json.Marshal(i.Products)
	if err != nil {
		return nil, err
	}

	sqlStatement := `
INSERT INTO orders (customer, status, products, total)
VALUES ($1, $2, $3, $4) RETURNING order_id`

	var orderId string
	err = cfg.DB.QueryRow(sqlStatement, customerJSON, "PENDING", productsJSON, i.Total).Scan(&orderId)
	if err != nil {
		return nil, err
	}
	return &orderId, nil
}
