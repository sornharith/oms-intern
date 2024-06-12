package service

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "New"
	OrderStatusPaid       OrderStatus = "Paid"
	OrderStatusProcessing OrderStatus = "Processing"
	OrderStatusDone       OrderStatus = "Done"
)
