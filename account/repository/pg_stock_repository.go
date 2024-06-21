package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"memrizr/account/entity"
	"strings"
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
	err := r.DB.Get(&stock, `SELECT s.s_id as SID, s.quantity as Quantity FROM stocks s 
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

func (r *stockRepository) DeductStockBulk(deductions map[int]int) error {
	if len(deductions) == 0 {
		return nil
	}

	var productIDs []string
	var amounts []string
	var params []interface{}
	var index = 1

	for productID, amount := range deductions {
		productIDs = append(productIDs, fmt.Sprintf("$%d", index))
		amounts = append(amounts, fmt.Sprintf("$%d", index+1))
		params = append(params, productID, amount)
		index += 2
	}

	// Step 1: Validate stock levels
	queryValidation := fmt.Sprintf(`
        WITH product_amounts AS (
            SELECT
                unnest(ARRAY[%s])::int AS product_id,
                unnest(ARRAY[%s])::int AS amount
        )
        SELECT
            p.p_id AS product_id,
            s.quantity,
            pa.amount
        FROM
            stocks s
            JOIN products p ON s.s_id = p.s_id
            JOIN product_amounts pa ON p.p_id = pa.product_id
        FOR UPDATE;`,
		strings.Join(productIDs, ", "), strings.Join(amounts, ", "))

	var results []struct {
		ProductID int `db:"product_id"`
		Quantity  int `db:"quantity"`
		Amount    int `db:"amount"`
	}

	err := r.DB.Select(&results, queryValidation, params...)
	if err != nil {
		return fmt.Errorf("error during validation phase: %w", err)
	}

	for _, result := range results {
		if result.Quantity < result.Amount {
			return fmt.Errorf("insufficient stock for product ID %d: available %d, required %d",
				result.ProductID, result.Quantity, result.Amount)
		}
	}

	// Step 2: Perform bulk stock deduction
	queryUpdate := fmt.Sprintf(`
        WITH product_amounts AS (
            SELECT
                unnest(ARRAY[%s])::int AS product_id,
                unnest(ARRAY[%s])::int AS amount
        ),
        updated AS (
            SELECT
                s.s_id,
                pa.product_id,
                pa.amount,
                s.quantity - pa.amount AS new_quantity
            FROM
                stocks s
                JOIN products p ON s.s_id = p.s_id
                JOIN product_amounts pa ON p.p_id = pa.product_id
            WHERE s.quantity - pa.amount >= 0
        )
        UPDATE stocks
        SET quantity = updated.new_quantity
        FROM updated
        WHERE stocks.s_id = updated.s_id
        RETURNING updated.s_id;`,
		strings.Join(productIDs, ", "), strings.Join(amounts, ", "))

	rows, err := r.DB.Queryx(queryUpdate, params...)
	if err != nil {
		return fmt.Errorf("error during bulk update phase: %w", err)
	}
	defer rows.Close()

	updatedProductIDs := make(map[int]bool)
	for rows.Next() {
		var sID int
		if err := rows.Scan(&sID); err != nil {
			return fmt.Errorf("error scanning updated rows: %w", err)
		}
		updatedProductIDs[sID] = true
	}

	if len(updatedProductIDs) != len(deductions) {
		return errors.New("insufficient stock for one or more products")
	}

	return nil
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
