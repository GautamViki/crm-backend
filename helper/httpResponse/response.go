package httpresponse

type ValidationErrors struct {
	ErrorMessages []string `json:"messages"`
}

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ValidationErrors
}

// Prepare response
func PrepareResponse(code string, message string) Response {
	res := Response{}
	res.Code = code
	if message != "" {
		res.Message = message
	}
	res.ErrorMessages = []string{}
	return res
}
