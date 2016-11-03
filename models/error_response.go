package models

type ErrorResponse struct {
	Code int `json:"code"`
}

func MakeError(code int) *ErrorResponse {
	return &ErrorResponse{
		Code: code,
	}
}
