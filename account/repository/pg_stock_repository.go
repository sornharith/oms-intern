package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"memrizr/account/entity"
)

type stockRepository struct {
	DB *sqlx.DB
}

func NewStockRepository(db *sqlx.DB) entity.StockRepository {
	return &stockRepository{
		DB: db,
	}
}

func (r *stockRepository) GetStockByProductID(productID int) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.DB.Get(&stock, `SELECT s.* FROM stocks s 
                             JOIN products p ON s.s_id = p.s_id 
                             WHERE p.p_id = $1`, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return &stock, nil
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

func (r *stockRepository) UpdateStock(StockId int, amount int) error {
	// Corrected SQL query with parameters in the right order
	query := `UPDATE stocks SET quantity = $1 WHERE s_id = $2`

	// Execute the query with the correct parameter order
	result, err := r.DB.Exec(query, amount, StockId)
	if err != nil {
		return fmt.Errorf("error updating stock: %w", err)
	}

	// Check if any rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no stock found with the given ID or no update needed")
	}

	return nil
}
