package helper

// Response codes
const (
	APISuccessCode                    = "0"
	EmailInvalidError                 = "Email is invalid."
	PhoneInvalidError                 = "Phone is invalid."
	FileRetrieveFromFormDataErrorCode = "1000"
	FileRetrieveFromFormDataError     = "No file was provided in the request. Please upload a file."
	FileFormateInvalidErrorCode       = "1001"
	FileFormateInvalidError           = "File format invalid error."
)

const (
	// code
	CustomerSaveErrorCode        = "2000"
	CustomerFetchErrorCode       = "2002"
	CustomerIdInvalidErrorCode   = "2003"
	CustomerDataInvalidErrorCode = "2004"
	CustomerUpdateErrorCode      = "2005"
	CustomerDeleteErrorCode      = "2006"
	CustomerNotExistErrorCode    = "2007"

	// message
	CustomerSaveError        = "Error occurred while saving customer."
	CustomerSaveSuccess      = "Save customers successfully."
	CustomerFetchError       = "Error occurred while fetching customers."
	CustomerFetchSuccess     = "Customers fetched successfully."
	CustomerIdInvalidError   = "Invalid customer id."
	CustomerDataInvalidError = "Invalid customer data."
	CustomerUpdateError      = "Error occurred while updating customer."
	CustomerUpdateSuccess    = "Customer updated successfully."
	CustomerDeleteError      = "Error occurred while deleting customer."
	CustomerNotExistError    = "Customer does not exist."
	CustomerDeleteSuccess    = "Customer deleted successfully."

	// code
	ExcelFileParseErrorCode = "3000"
	// message
	ExcelFileOpenError            = "Error occurred while opening excel file."
	ExcelFileRowReadError         = "Error occurred while reading row of excel file."
	ExcelFileEmptyError           = "Excel file is empty."
	ExcelColumnHeaderInvalidError = "Invalid column header."
	ExcelCulumnInsufficientError  = "Row has insufficient headers."
	ExcelFileParseError           = "Failed to parse Excel file."
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

const (
	XlsxFormat    = ".xlsx"
	DirectoryName = "upload"
)
