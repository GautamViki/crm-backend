package controller

import (
	"crm-backend/config"
	h "crm-backend/helper"
	httpresponse "crm-backend/helper/httpResponse"
	"crm-backend/models"
	"crm-backend/services"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"net/http"
)

type customerHandler struct {
	cfg *config.Config
}

func NewHandler() *customerHandler {
	return &customerHandler{cfg: &config.Config{}}
}

// ListCustomers fetches customers, first checking Redis
func (ch *customerHandler) ListCustomers(c *gin.Context) {
	config := config.ConnectDB()
	db := config.DB

	customers, err := services.GetAllCustomersFromCache(db, config.Redis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (ch *customerHandler) CreateCustomer(c *gin.Context) {
	config := config.ConnectDB()
	file, err := c.FormFile("file")
	if err != nil {
		res := httpresponse.PrepareResponse(h.FileFormateInvalidError, h.FileRetrieveFromFormDataError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	// Validate file type (must be .xlsx)
	if !strings.HasSuffix(file.Filename, h.XlsxFormat) {
		res := httpresponse.PrepareResponse(h.FileFormateInvalidErrorCode, h.FileFormateInvalidError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	customers, err := services.ParseExcel(file)
	if err != nil {
		res := httpresponse.PrepareResponse(h.ExcelFileParseErrorCode, h.ExcelFileParseError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(customers)) // Channel to collect errors

	for _, customer := range customers {
		wg.Add(1) // Increment the wait group counter
		go func(cust models.Customer, wg *sync.WaitGroup) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			if err := services.AddCustomer(&cust, config); err != nil {
				errCh <- err // Send the error to the channel
				return
			}
		}(customer, &wg)
	}

	wg.Wait()    // Wait for all goroutines to finish
	close(errCh) // Close the error channel

	// Check for errors
	if err := <-errCh; err != nil {
		res := httpresponse.PrepareResponse(h.CustomerSaveErrorCode, h.CustomerSaveError)
		h.RespondWithError(c, http.StatusInternalServerError, res)
		return
	}

	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerSaveSuccess)
	h.RespondWithError(c, http.StatusCreated, res)
}

// EditCustomer updates a customer in MySQL and Redis
func (ch *customerHandler) EditCustomer(c *gin.Context) {
	config := config.ConnectDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Update the customer in MySQL
	err = services.UpdateCustomer(id, &customer, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// EditCustomer updates a customer in MySQL and Redis
func (ch *customerHandler) DeleteCustomer(c *gin.Context) {
	config := config.ConnectDB()
	id, _ := strconv.Atoi(c.Param("id"))
	var customer models.Customer

	err := services.DeleteCustomer(id, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// GetAllCacheCustomers fetches all customers from Redis or MySQL
func (ch *customerHandler) GetAllCacheCustomers(c *gin.Context) {
	config := config.ConnectDB()
	db := config.DB
	redisClient := config.Redis // Assuming you have a Redis client configured

	customers, err := services.GetAllCustomersFromCache(db, redisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, customers)
}
