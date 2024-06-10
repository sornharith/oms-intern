package repository

import (
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
	_, err := r.DB.Exec(`UPDATE stocks SET quantity = quantity - $1 
                         WHERE s_id = (SELECT s_id FROM products WHERE p_id = $2)`, amount, productID)
	return err
}
