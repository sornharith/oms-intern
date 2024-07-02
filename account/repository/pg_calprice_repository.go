package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"memrizr/account/entity"
	"memrizr/account/entity/apperrors"
)

type calPriceRepository struct {
	DB *sqlx.DB
}

func NewCalPriceRepository(db *sqlx.DB) *calPriceRepository {
	return &calPriceRepository{
		DB: db,
	}
}

func (r *calPriceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	var calPrice entity.CalPrice
	err := r.DB.Get(&calPrice, "SELECT t_id as TID,t_price as TPrice, user_select as UserSelect, address as Address FROM calprice WHERE t_id::text=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("the calculate price not found with the provided ID")
		}
		return nil, fmt.Errorf("error on querying calprice: %w", err)
	}

	return &calPrice, nil
}

func (r *calPriceRepository) Update(ctx context.Context, calPrice *entity.CalPrice) error {
	query := "UPDATE calprice SET t_price=$1, user_select=$2 WHERE t_id=$3"
	_, err := r.DB.Exec(query, calPrice.TPrice, calPrice.UserSelect, calPrice.TID)

	return err
}

func (r *calPriceRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM calprice WHERE t_id=$1"
	_, err := r.DB.Exec(query, id)
	return err
}
func (r *calPriceRepository) CalculateTotalPrice(ctx context.Context, userSelect []map[string]interface{}) (float64, error) {
	var totalPrice float64
	var price float64
	for _, item := range userSelect {
		productID := int(item["product_id"].(float64))
		amount := int(item["amount"].(float64))

		if err := r.DB.Get(&price, "SELECT p_price FROM products WHERE p_id = $1", productID); err != nil {
			return 0, errors.New("product not found")
		}

		totalPrice += price * float64(amount)
	}

	return totalPrice, nil
}

func (r *calPriceRepository) CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) error {
	userSelectJSON, err := json.Marshal(calPrice.UserSelect)
	if err != nil {
		return err
	}
	query := "INSERT INTO calprice (t_price, user_select,address) VALUES ($1, $2, $3) RETURNING t_id"
	if err := r.DB.QueryRow(query, calPrice.TPrice, userSelectJSON, calPrice.Address).Scan(&calPrice.TID); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with id: %v. Reason: %v\n", calPrice.TID, err.Code.Name())
			return apperror.NewConflict("transaction is", "duplicate")
		}

		log.Printf("Could not create a user with id: %v. Reason: %v\n", calPrice.TID, err)
		return apperror.NewInternal()
	}
	return nil

}
