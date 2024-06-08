package model

type Stock struct {
	SID      int `db:"s_id" json:"s_id"`
	Quantity int `db:"quantity" json:"quantity"`
}
