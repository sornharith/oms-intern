package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"memrizr/account/entity"
	"net/http"
	"strconv"
)

// Handler struct holds required services for handler to function
type Handler struct {
	UserService     entity.UserService
	CalpriceService entity.CalPriceService
	ProductService  entity.ProductService
	OrderService    entity.OrderService
	StockService    entity.StockService
	Logger          *logrus.Logger
	Tracer          trace.Tracer
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     entity.UserService
	CalpriceService entity.CalPriceService
	ProductService  entity.ProductService
	OrderService    entity.OrderService
	StockService    entity.StockService
	Logger          *logrus.Logger
	Tracer          trace.Tracer
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
		Logger:          c.Logger,
		Tracer:          c.Tracer,
	}

	// Create an account group
	g := c.R
	//g.Use(logger.GinLogger(h.Logger), logger.GinRecovery(h.Logger), tracing.GinTracer())

	g.GET("/me", h.Me)
	// for the stock
	g.GET("/stock/:p_id", h.GetStockbyid)
	g.PATCH("/stock/:p_id", h.updateStock)
	// for prduct
	g.GET("/products", h.getproduct)
	// for the calprice
	g.GET("/transaction/:t_id", h.GettransactionbyID)
	g.POST("/order/calculate", h.CreateCalPrice)
	// for the order
	g.POST("/order", h.CreateOrder)
	g.PATCH("/order/status/:o_id", h.UpdateOrderStatus)
	g.GET("/order/:o_id", h.GetorderById)

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
	ctx := c.Request.Context()
	calp, err := h.CalpriceService.GetCalPriceByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Transaction information": calp,
	})
}

func (h *Handler) GetStockbyid(c *gin.Context) {
	type Stockoutput struct {
		Product_id int `json:"product_id"`
		Stock_id   int `json:"stock_id"`
		Quantity   int `json:"quantity"`
	}
	pid := c.Param("p_id")
	id, err := strconv.Atoi(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productid"})
		return
	}
	ctx := c.Request.Context()
	stock, err := h.StockService.GetStockByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stockresponse := Stockoutput{
		Product_id: id,
		Stock_id:   stock.SID,
		Quantity:   stock.Quantity,
	}
	c.JSON(http.StatusOK, gin.H{
		"Stock information": stockresponse,
	})
}

func (h *Handler) getproduct(c *gin.Context) {
	product, err := h.ProductService.GetallProductwithstock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": product,
	})
}

func (h *Handler) GetorderById(c *gin.Context) {
	oid := c.Param("o_id")
	id, err := uuid.Parse(oid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx := c.Request.Context()
	order, err := h.OrderService.GetOrderByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"errorid": "400",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Order information": order,
	})
}
