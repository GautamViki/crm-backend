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
	DirectoryCreateErrorCode          = "1002"
	DirectoryCreateError              = "Unable to create directory"
	FileSaveErrorCode                 = "1003"
	FileSaveError                     = "Failed to save file."
	ValidateExcelFileErrorCode        = "1004"
	FileUploadAndValidateSuccess      = "File uploaded and validated successfully."
)

const (
	// code
	ExcelFileParseErrorCode = "2000"
	// message
	ExcelFileOpenError            = "Error occurred while opening excel file."
	ExcelFileRowReadError         = "Error occurred while reading row of excel file."
	ExcelFileEmptyError           = "Excel file is empty."
	ExcelColumnHeaderInvalidError = "Invalid column header."
	ExcelCulumnInsufficientError  = "Row has insufficient headers."
	ExcelFileParseError           = "Failed to parse Excel file."

	CustomerSaveErrorCode = "2000"
	CustomerSaveError     = "Error occurred while saving customer."
	CustomerSaveSuccess="Save customers successfully."
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
