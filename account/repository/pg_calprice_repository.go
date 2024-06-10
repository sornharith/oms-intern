package repository

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"memrizr/account/model"
	apperror "memrizr/account/model/apperrors"
)

type calPriceRepository struct {
	DB *sqlx.DB
}

func NewCalPriceRepository(db *sqlx.DB) *calPriceRepository {
	return &calPriceRepository{
		DB: db,
	}
}

func (r *calPriceRepository) GetByID(id uuid.UUID) (*model.CalPrice, error) {
	var calPrice model.CalPrice
	err := r.DB.Get(&calPrice, "SELECT * FROM calprice WHERE t_id::text=$1", id)
	if err != nil {
		log.Printf("error on querying calprice %v", err.Error())
		return nil, err
	}
	return &calPrice, nil
}

func (r *calPriceRepository) Update(calPrice *model.CalPrice) error {
	query := "UPDATE calprice SET t_price=$1, user_select=$2 WHERE t_id=$3"
	_, err := r.DB.Exec(query, calPrice.TPrice, calPrice.UserSelect, calPrice.TID)

	return err
}

func (r *calPriceRepository) Delete(id int) error {
	query := "DELETE FROM calprice WHERE t_id=$1"
	_, err := r.DB.Exec(query, id)
	return err
}
func (r *calPriceRepository) CalculateTotalPrice(userSelect []map[string]interface{}) (float64, error) {
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

func (r *calPriceRepository) CreateCalPrice(calPrice *model.CalPrice) error {
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
