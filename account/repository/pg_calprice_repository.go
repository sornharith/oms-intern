package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memrizr/account/entity"
	"memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type calPriceRepository struct {
	DB *sqlx.DB
}

func NewCalPriceRepository(db *sqlx.DB) CalPriceRepository {
	return &calPriceRepository{
		DB: db,
	}
}

func (r *calPriceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	_, span := tracer.Start(ctx, "repository get-calprice-by-id")
	defer span.End()
	var calPrice entity.CalPrice
	err := r.DB.Get(&calPrice, "SELECT t_id as TID,t_price as TPrice, user_select as UserSelect, address as Address FROM calprice WHERE t_id::text=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.LogError(apperror.CusNotFound("not found calculate price id "+id.String(), "3044"), "not found calculate price", logrus.Fields{
				"at": "repository",
			})
			return nil, errors.New("the calculate price not found with the provided ID")
		}
		logger.LogError(apperror.CusBadRequest("error on querying", "3140"), "error on querying", logrus.Fields{
			"at": "repository",
		})	
		return nil, fmt.Errorf("error on querying calprice: %w", err)
	}

	return &calPrice, nil
}

func (r *calPriceRepository) Update(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error) {
	query := "UPDATE calprice SET t_price=$1, user_select=$2 WHERE t_id=$3"
	_, err := r.DB.Exec(query, calPrice.TPrice, calPrice.UserSelect, calPrice.TID)
	if err != nil {
		logger.LogError(apperror.CusBadRequest("can't update calculate price", "3040"), "can't update calculate price", logrus.Fields{
			"at": "repository",	
		})
	}
	return calPrice, err
}

func (r *calPriceRepository) Delete(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	res, _ := r.GetByID(ctx, id)
	query := "DELETE FROM calprice WHERE t_id::text=$1"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		logger.LogError(apperror.CusBadRequest("can't delete calculate price", "3040"), "can't delete calculate price", logrus.Fields{
			"at": "repository",
		})
	}
	return res, err
}
func (r *calPriceRepository) CalculateTotalPrice(ctx context.Context, userSelect []map[string]interface{}) (float64, error) {
	_, span := tracer.Start(ctx, "repository calculate-total-price")
	defer span.End()

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

func (r *calPriceRepository) CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error) {
	_, span := tracer.Start(ctx, "repository create-calprice")
	defer span.End()
	userSelectJSON, err := json.Marshal(calPrice.UserSelect)
	if err != nil {
		return nil, err
	}
	query := "INSERT INTO calprice (t_price, user_select,address) VALUES ($1, $2, $3) RETURNING t_id"
	if err := r.DB.QueryRow(query, calPrice.TPrice, userSelectJSON, calPrice.Address).Scan(&calPrice.TID); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			logger.LogError(apperror.CusBadRequest("transaction is duplicate", "3040"), "transaction is duplicate", logrus.Fields{
				"at": "repository",
			})
			log.Printf("Could not create a calprice with id: %v. Reason: %v\n", calPrice.TID, err.Code.Name())
			return nil, apperror.NewConflict("transaction is", "duplicate")
		}

		logger.LogError(apperror.CusBadRequest("can't create calculate price", "3140"), "can't create calculate price", logrus.Fields{
			"at": "repository",
		})
		log.Printf("Could not create a calprice with id: %v. Reason: %v\n", calPrice.TID, err)
		return nil, apperror.NewInternal()
	}
	return calPrice, nil

}
