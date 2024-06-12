package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"memrizr/account/entity"
	"net/http"
	"strconv"
)

// Handler struct holds required services for handler to function
type Handler struct {
	UserService     entity.UserService
	CalpriceService entity.CalPriceService
	ProductService  entity.ProductRepository
	OrderService    entity.OrderService
	StockService    entity.StockService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     entity.UserService
	CalpriceService entity.CalPriceService
	ProductService  entity.ProductRepository
	OrderService    entity.OrderService
	StockService    entity.StockService
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
		StockService:    c.StockService,
	}

	// Create an account group
	g := c.R

	g.GET("/me", h.Me)
	// for the stock
	g.GET("/stock/:p_id", h.GetStockbyid)
	// for the calprice
	g.GET("/transaction/:t_id", h.GettransactionbyID)
	g.POST("/order/calculate", h.CreateCalPrice)
	// for the order
	g.POST("/order", h.CreateOrder)
	g.PATCH("/order/status/:o_id", h.UpdateOrderStatus)

	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

}

func (h *Handler) GettransactionbyID(c *gin.Context) {
	tid := c.Param("t_id")
	id, err := uuid.Parse(tid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	calp, err := h.CalpriceService.GetCalPriceByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Transaction information": calp,
	})
}

func (h *Handler) GetStockbyid(c *gin.Context) {
	pid := c.Param("p_id")
	id, err := strconv.Atoi(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productid"})
		return
	}
	stock, err := h.StockService.GetStockByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Stock information": stock,
	})

}
