package dto

import (
	httpresponse "crm-backend/helper/httpResponse"
	"crm-backend/models"
)

type CustomersResponse struct {
	httpresponse.Response
	Total     int               `json:"total"`
	Customers []models.Customer `json:"customers"`
}
type CustomerResponse struct {
	httpresponse.Response
	Customer models.Customer `json:"customer"`
}
