package handler

import (
	"github.com/gin-gonic/gin"
	"memrizr/account/model"
)

// Handler struct holds required services for handler to function
type Handler struct {
	UserService     model.UserService
	CalpriceService model.CalPriceService
	ProductService  model.ProductRepository
	OrderService    model.OrderService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	CalpriceService model.CalPriceService
	ProductService  model.ProductRepository
	OrderService    model.OrderService
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	// Create a handler (which will later have injected services)
	h := &Handler{
		UserService:     c.UserService,
		CalpriceService: c.CalpriceService,
		ProductService:  c.ProductService,
		OrderService:    c.OrderService,
	}

	// Create an account group
	g := c.R

	g.GET("/me", h.Me)
	//g.POST("/signup", h.Signup)
	g.POST("/order/calculate", h.CreateCalPrice)
	g.POST("/order", h.CreateOrder)
	g.PATCH("/order/status/:o_id", h.UpdateOrderStatus)
}
