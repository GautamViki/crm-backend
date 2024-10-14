package main

import (
	"crm-backend/config"
	"crm-backend/controller"
	"crm-backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.ConnectDB()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())
	// Initialize database and Redis connections
	// cfg := c

	// Migrate the schema
	models.Migrate(config.DB)

	// Initialize Gin router
	controllers := controller.NewHandler()
	// Routes for CRUD
	r.POST("/customers", controllers.CreateCustomer)
	r.GET("/customers", controllers.ListCustomers)
	r.PUT("/customers/:id", controllers.EditCustomer)
	r.DELETE("/customers/:id", controllers.DeleteCustomer)
	r.GET("/customers/cache", controllers.GetAllCacheCustomers)

	// Start the server
	r.Run(":8080")
}
