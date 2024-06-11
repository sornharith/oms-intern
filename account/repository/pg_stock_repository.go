package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"memrizr/account/model"
)

type stockRepository struct {
	DB *sqlx.DB
}

func NewStockRepository(db *sqlx.DB) model.StockRepository {
	return &stockRepository{
		DB: db,
	}
}

func (r *stockRepository) GetStockByProductID(productID int) (*model.Stock, error) {
	var stock model.Stock
	err := r.DB.Get(&stock, `SELECT s.* FROM stocks s 
                             JOIN products p ON s.s_id = p.s_id 
                             WHERE p.p_id = $1`, productID)
	return &stock, err
}

func (r *stockRepository) DeductStock(productID int, amount int) error {
	query := `UPDATE stocks
              SET quantity = quantity - $1
              WHERE s_id = (SELECT s_id FROM products WHERE p_id = $2)
              AND quantity >= $1`

	res, err := r.DB.Exec(query, amount, productID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}

func (r *stockRepository) AddStock(productID int, amount int) error {
	query := `UPDATE stocks
              SET quantity = quantity + $1
              WHERE s_id = (SELECT s_id FROM products WHERE p_id = $2)`

	_, err := r.DB.Exec(query, amount, productID)
	return err
}
