package service

type ProductStock struct {
	ProductId int `json:"product_id"`
	StockId   int `json:"stock_id"`
	Quantity  int `json:"quantity"`
}
