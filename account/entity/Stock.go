package entity

type Stock struct {
	SID      int `db:"s_id"`
	Quantity int `db:"quantity"`
}
