package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"memrizr/account/entity"
)

type orderRepository struct {
	DB *sqlx.DB
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.DB.Get(&order, "SELECT o_id as OID, t_id as TID, t_price as TPrice, status as Status, create_at as CreatedAt, last_edit as LastEdit FROM orders WHERE o_id::text = $1", id)
	if err != nil {
		log.Printf("error on querying calprice %v", err.Error())
		return nil, err
	}
	return &order, err
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	query := "UPDATE orders SET status = $1, last_edit = CURRENT_TIMESTAMP WHERE o_id::text = $2;"
	_, err := r.DB.Exec(query, order.Status, order.OID)
	return err
}

func (r *orderRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	query := `INSERT INTO orders (t_id, t_price, status, create_at, last_edit) 
              VALUES ($1, $2, $3, $4, $5) RETURNING o_id`
	return r.DB.QueryRow(query, order.TID, order.TPrice, order.Status, order.CreatedAt, order.LastEdit).Scan(&order.OID)
}
