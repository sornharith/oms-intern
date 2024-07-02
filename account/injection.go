package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"log"
	"memrizr/account/handler"
	"memrizr/account/logger"
	"memrizr/account/repository"
	"memrizr/account/service"
	"memrizr/account/tracing"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(d.DB)
	calPriceRepo := repository.NewCalPriceRepository(d.DB)
	productRepo := repository.NewProductRepository(d.DB)
	orderRepo := repository.NewOrderRepository(d.DB)
	stockRepo := repository.NewStockRepository(d.DB)
	/*
	 * repository layer
	 */
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	calPriceService := service.NewCalPriceUsecase(&service.CalpConfig{
		CalPriceRepo: calPriceRepo,
		ProductRepo:  productRepo,
	})

	orderService := service.NewCreateOrderUsecase(&service.CreateOrderconfig{
		CalPriceRepo: calPriceRepo,
		OrderRepo:    orderRepo,
		StockRepo:    stockRepo,
	})

	stockService := service.NewStockService(&service.StockConfig{
		StockRepository: stockRepo,
	})

	productService := service.NewProductService(&service.ProductConfig{
		ProductRepository: productRepo,
	})

	logger.Setup()
	tracerShutdown := tracing.InitTracer()
	defer tracerShutdown()

	// initialize gin.Engine
	router := gin.Default()

	//router.Use(logger.GinLogger(log), logger.GinRecovery(log), tracing.GinTracer())

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		CalpriceService: calPriceService,
		ProductService:  productService,
		OrderService:    orderService,
		StockService:    stockService,
		Tracer:          otel.Tracer("oms-handler"),
	})

	return router, nil
}
