package repository

import (
	"context"
	"errors"
	"fmt"
	"memrizr/account/entity"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type stockRepository struct {
	DB *sqlx.DB
}

func NewStockRepository(db *sqlx.DB) StockRepository {
	return &stockRepository{
		DB: db,
	}
}

func (r *stockRepository) GetStockByProductID(ctx context.Context, productID int) (*entity.Stock, error) {
	_, span := tracer.Start(ctx, "repository get-stock-by-product-id")
	defer span.End()

	var stock entity.Stock
	err := r.DB.Get(&stock, `SELECT s.s_id as SID, s.quantity as Quantity FROM stocks s 
                             JOIN products p ON s.s_id = p.s_id 
                             WHERE p.p_id = $1`, productID)
	if err != nil {
		logger.LogError(apperror.CusNotFound(strconv.Itoa(productID), "3044"), "not found from this product", logrus.Fields{
			"at": "repository",
		})
		return nil, errors.New("product not found")
	}
	return &stock, nil
}

func (r *stockRepository) DeductStock(ctx context.Context, productID int, amount int) error {
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

//func (r *stockRepository) AddStock(ctx context.Context, productID int, amount int) error {
//	query := `UPDATE stocks
//              SET quantity = quantity + $1
//              WHERE s_id = (SELECT s_id FROM products WHERE p_id = $2)`
//
//	_, err := r.DB.Exec(query, amount, productID)
//	return err
//}

func (r *stockRepository) DeductStockBulk(ctx context.Context, deductions map[int]int) error {
	_, span := tracer.Start(ctx, "repository deduct-stock-bulk")
	defer span.End()

	if len(deductions) == 0 {
		logger.LogError(apperror.CusBadRequest("no deduction request", "3040"), "no deduction request", logrus.Fields{
			"at": "repository",
		})
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
		logger.LogError(apperror.CusNotFound(err.Error(), "3144"), "error during validation phase", logrus.Fields{
			"at": "repository",
		})
		return fmt.Errorf("error during validation phase: %w", err)
	}

	for _, result := range results {
		if result.Quantity < result.Amount {
			logger.LogError(apperror.CusBadRequest("insufficient stock", "3240"), "error during validation phase", logrus.Fields{
				"at": "repository",
			})
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
		logger.LogError(apperror.CusBadRequest(err.Error(), "3340"), "error during bulk update phase", logrus.Fields{
			"at": "repository",
		})
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
		logger.LogError(apperror.CusBadRequest("insufficient stock", "3440"), "error during bulk update phase", logrus.Fields{
			"at": "repository",
		})
		return errors.New("insufficient stock for one or more products")
	}

	return nil
}

func (r *stockRepository) UpdateStock(ctx context.Context, stock *entity.Stock) (*entity.Stock, error) {
	_, span := tracer.Start(ctx, "repository update-stock")
	defer span.End()

	// Corrected SQL query with parameters in the right order
	query := `UPDATE stocks SET quantity = $1 WHERE s_id = $2`

	// Execute the query with the correct parameter order
	result, err := r.DB.Exec(query, stock.Quantity, stock.SID)
	if err != nil {
		return nil, fmt.Errorf("error updating stock: %w", err)
	}

	// Check if any rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, errors.New("no stock found with the given ID or no update needed")
	}

	return stock, nil
}
