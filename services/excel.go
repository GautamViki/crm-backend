package services

import (
	h "crm-backend/helper"
	"crm-backend/models"
	"errors"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

var ExcelFileHeader = []string{
	"first_name",
	"last_name",
	"company_name",
	"address",
	"city",
	"county",
	"postal",
	"phone",
	"email",
	"web",
}

func ParseExcel(file *multipart.FileHeader) ([]models.Customer, error) {
	f, err := file.Open()
	if err != nil {
		return nil, errors.New(h.ExcelFileOpenError)
	}
	defer f.Close()

	excel, err := excelize.OpenReader(f)
	if err != nil {
		return nil, errors.New(h.ExcelFileOpenError)
	}

	// Assuming the first sheet is used
	sheetName := excel.GetSheetList() // Get the name of the first sheet
	rows, err := excel.GetRows(sheetName[0])
	if err != nil {
		return []models.Customer{}, errors.New(h.ExcelFileRowReadError)
	}
	if len(rows) == 0 {
		return []models.Customer{}, errors.New(h.ExcelFileEmptyError)
	}

	// Check column headers (assuming we expect these exact headers)
	for i, header := range ExcelFileHeader {
		if i >= len(rows[0]) || rows[0][i] != header {
			return []models.Customer{}, errors.New(h.ExcelColumnHeaderInvalidError)
		}
	}
	customers := []models.Customer{}
	for _, row := range rows[1:] {
		customers = append(customers, models.Customer{
			FirstName: row[0], LastName: row[1], Company: row[2],
			Address: row[3], City: row[4], County: row[5],
			Postal: row[6], Phone: row[7], Email: row[8], Web: row[9],
		})
	}

	return customers, nil
}
