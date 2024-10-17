package services

import (
	h "crm-backend/helper"
	"crm-backend/models"
	"errors"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

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

	sheetName := excel.GetSheetList()
	rows, err := excel.GetRows(sheetName[0])
	if err != nil {
		return []models.Customer{}, errors.New(h.ExcelFileRowReadError)
	}
	if len(rows) == 0 {
		return []models.Customer{}, errors.New(h.ExcelFileEmptyError)
	}

	if len(rows[0]) != len(h.ExcelFileHeader) {
		return []models.Customer{}, errors.New(h.ExcelCulumnInsufficientError)
	}

	for i, header := range h.ExcelFileHeader {
		if i >= len(rows[0]) || rows[0][i] != header {
			return []models.Customer{}, errors.New(h.ExcelColumnHeaderInvalidError)
		}
	}

	customers := []models.Customer{}
	for _, row := range rows[1:] {
		if !h.ValidateEmail(row[8]) {
			return []models.Customer{}, errors.New(h.EmailInvalidError)
		}
		// if !h.ValidatePhoneNumberWithCountry(row[7]) {
		// 	return []models.Customer{}, errors.New(h.PhoneInvalidError)
		// }
		customers = append(customers, models.Customer{
			FirstName: row[0], LastName: row[1], Company: row[2],
			Address: row[3], City: row[4], County: row[5],
			Postal: row[6], Phone: row[7], Email: row[8], Web: row[9],
		})
	}
	return customers, nil
}
