package entity

type Product struct {
	PID    int     `db:"p_id" json:"p_id"`
	PName  string  `db:"p_name" json:"p_name"`
	PDesc  string  `db:"p_desc" json:"p_desc"`
	PPrice float64 `db:"p_price" json:"p_price"`
	SID    int     `db:"s_id" json:"s_id"`
}
