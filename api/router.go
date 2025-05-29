package api

import (
	"net/http"
	"parte3/internal/sale"
	"parte3/internal/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InitRoutes registers all user CRUD endpoints on the given Gin engine.
// It initializes the storage, service, and handler, then binds each HTTP
// method and path to the appropriate handler function.
func InitRoutes(e *gin.Engine) {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	storage := user.NewLocalStorage()
	service := user.NewService(storage, logger)
	salesStorage := sale.NewLocalStorage()
	salesService := sale.NewService(salesStorage, service, logger)
	// Initialize handler with services
	h := handler{
		userService: service,
		logger:      logger,
		saleService: salesService,
	}

	e.POST("/users", h.handleCreate)
	e.POST("/sales", h.handleCreateSale)
	e.GET("/users/:id", h.handleRead)
	e.GET("/users", h.handleListActive)
	e.GET("/sales/:id/:status", h.handleReadSales)
	e.PATCH("/users/:id", h.handleUpdate)
	e.DELETE("/users/:id", h.handleDelete)
	e.PATCH("/sales/:id", h.handleUpdateSaleStatus)

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
