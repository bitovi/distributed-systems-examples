package activities

import (
	"context"
	"database/sql"
)

type UpdateOrderConfig struct {
	DB *sql.DB
}

func NewUpdateOrder(db *sql.DB) *UpdateOrderConfig {
	return &UpdateOrderConfig{
		DB: db,
	}
}

func (cfg UpdateOrderConfig) UpdateOrderActivity(ctx context.Context, status, oid string) (*string, error) {
	sqlStatement := `
	UPDATE orders SET status = $1
	where order_id = $2`

	_, err := cfg.DB.Exec(sqlStatement, status, oid)
	if err != nil {
		return nil, err
	}
	return &status, nil
}
