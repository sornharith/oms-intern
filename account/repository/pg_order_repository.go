package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"memrizr/account/model"
)

type orderRepository struct {
	DB *sqlx.DB
}

func (r *orderRepository) GetByID(id int) (*model.Order, error) {
	var order model.Order
	err := r.DB.Get(&order, "SELECT * FROM orders WHERE o_id = $1", id)
	if err != nil {
		log.Printf("error on querying calprice %v", err.Error())
		return nil, err
	}
	return &order, err
}

func (r *orderRepository) Create(order *model.Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *orderRepository) Update(order *model.Order) error {
	query := "UPDATE orders SET status = $1, last_edit = CURRENT_TIMESTAMP WHERE o_id = $2;"
	_, err := r.DB.Exec(query, order.Status, order.OID)
	return err
}

func (r *orderRepository) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) CreateOrder(order *model.Order) error {
	query := `INSERT INTO orders (o_id, t_id, t_price, status, create_at, last_edit) 
              VALUES ($1, $2, $3, $4, $5,$6) RETURNING o_id`
	return r.DB.QueryRow(query, order.OID, order.TID, order.TPrice, order.Status, order.CreatedAt, order.LastEdit).Scan(&order.OID)
}
