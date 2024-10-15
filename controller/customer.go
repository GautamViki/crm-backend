package controller

import (
	"crm-backend/config"
	h "crm-backend/helper"
	httpresponse "crm-backend/helper/httpResponse"
	"crm-backend/internals/dto"
	"crm-backend/models"
	"crm-backend/services"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"net/http"
)

type customerHandler struct {
}

func NewHandler() *customerHandler {
	return &customerHandler{}
}

func (ch *customerHandler) ListCustomers(c *gin.Context) {
	dbConfig := config.ConnectDB()
	db := dbConfig.DB
	logger := config.GetLoggerInstance()
	customers, err := services.GetAllCustomersFromCache(db, dbConfig.Redis)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerFetchErrorCode, h.CustomerFetchError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	logger.Log(h.CustomerFetchSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerFetchSuccess)
	custResp := dto.CustomersResponse{
		Response:  res,
		Total:     len(customers),
		Customers: customers,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

func (ch *customerHandler) UploadCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	file, err := c.FormFile("file")
	if err != nil {
		logger.Log(h.FileRetrieveFromFormDataError, err.Error(), h.FileRetrieveFromFormDataErrorCode)
		res := httpresponse.PrepareResponse(h.FileRetrieveFromFormDataErrorCode, h.FileRetrieveFromFormDataError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	if !strings.HasSuffix(file.Filename, h.XlsxFormat) {
		logger.Log(h.FileFormateInvalidError, err.Error(), h.FileFormateInvalidErrorCode)
		res := httpresponse.PrepareResponse(h.FileFormateInvalidErrorCode, h.FileFormateInvalidError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	customers, err := services.ParseExcel(file)
	if err != nil {
		logger.Log(h.ExcelFileParseError, err.Error(), h.ExcelFileParseErrorCode)
		res := httpresponse.PrepareResponse(h.ExcelFileParseErrorCode, h.ExcelFileParseError)
		h.RespondWithError(c, http.StatusInternalServerError, res)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(customers))

	for _, customer := range customers {
		wg.Add(1)
		go func(cust models.Customer, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := services.AddCustomer(&cust, dbConfig); err != nil {
				errCh <- err
				return
			}
		}(customer, &wg)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		logger.Log(h.CustomerSaveError, err.Error(), h.CustomerSaveErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerSaveErrorCode, h.CustomerSaveError)
		h.RespondWithError(c, http.StatusInternalServerError, res)
		return
	}
	logger.Log(h.CustomerSaveSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerSaveSuccess)
	h.RespondWithJSON(c, res, http.StatusCreated)
}

func (ch *customerHandler) UpdateCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log(h.CustomerIdInvalidError, err.Error(), h.CustomerIdInvalidErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerIdInvalidErrorCode, h.CustomerIdInvalidError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}
	customer, err := services.FetchById(id, dbConfig.DB)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerFetchErrorCode, h.CustomerFetchError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		logger.Log(h.CustomerDataInvalidError, err.Error(), h.CustomerDataInvalidErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerDataInvalidErrorCode, h.CustomerDataInvalidError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	customer, err = services.UpdateCustomer(id, customer, dbConfig)
	if err != nil {
		logger.Log(h.CustomerUpdateError, err.Error(), h.CustomerUpdateErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerUpdateErrorCode, h.CustomerUpdateError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	logger.Log(h.CustomerUpdateSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerUpdateSuccess)
	custResp := dto.CustomerResponse{
		Response: res,
		Customer: customer,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

func (ch *customerHandler) DeleteCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log(h.CustomerIdInvalidError, err.Error(), h.CustomerIdInvalidErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerIdInvalidErrorCode, h.CustomerIdInvalidError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	customer, err := services.FetchById(id, dbConfig.DB)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerFetchErrorCode, h.CustomerFetchError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	err = services.DeleteCustomer(id, dbConfig)
	if err != nil {
		logger.Log(h.CustomerDeleteError, err.Error(), h.CustomerDeleteErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerDeleteErrorCode, h.CustomerDeleteError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	logger.Log(h.CustomerDeleteSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerDeleteSuccess)
	custResp := dto.CustomerResponse{
		Response: res,
		Customer: customer,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

func (ch *customerHandler) GetAllCacheCustomers(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	db := dbConfig.DB
	redisClient := dbConfig.Redis

	customers, err := services.GetAllCustomersFromCache(db, redisClient)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		res := httpresponse.PrepareResponse(h.CustomerFetchErrorCode, h.CustomerFetchError)
		h.RespondWithError(c, http.StatusBadRequest, res)
		return
	}

	logger.Log(h.CustomerFetchSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerFetchSuccess)
	custResp := dto.CustomersResponse{
		Response:  res,
		Total:     len(customers),
		Customers: customers,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}
