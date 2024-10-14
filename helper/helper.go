package helper

import (
	httpresponse "crm-backend/helper/httpResponse"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// RespondWithJSON sends a JSON response
func RespondWithJSON(c *gin.Context, res interface{}, statusCode int) {
	c.JSON(statusCode, res)
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, res httpresponse.Response) {
	RespondWithJSON(c, res, statusCode)
}

func ValidateExcelFile(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return errors.New(ExcelFileOpenError)
	}

	// Assuming the first sheet is used
	sheetName := f.GetSheetList() // Get the name of the first sheet
	rows, err := f.GetRows(sheetName[0])
	if err != nil {
		return errors.New(ExcelFileRowReadError)
	}

	if len(rows) == 0 {
		return errors.New(ExcelFileEmptyError)
	}

	// Check column headers (assuming we expect these exact headers)
	for i, header := range ExcelFileHeader {
		if i >= len(rows[0]) || rows[0][i] != header {
			return errors.New(ExcelColumnHeaderInvalidError)
		}
	}

	// Validate data rows
	for _, row := range rows[1:] {
		if len(row) < len(ExcelFileHeader) {
			return errors.New(ExcelCulumnInsufficientError)
		}

		if !ValidatePhoneNumberWithCountry(row[7]) {
			return errors.New(PhoneInvalidError)
		}

		if !ValidateEmail(row[8]) {
			return errors.New(EmailInvalidError)
		}
	}
	return nil
}

func ValidateEmail(email string) bool {
	// Regular expression for validating an email
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePhoneNumberWithCountry validates a phone number with an optional country code
func ValidatePhoneNumberWithCountry(phone string) bool {
	const phoneRegex = `^[+]{1}(?:[0-9\-\(\)\/\.]\s?){6, 15}[0-9]{1}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
