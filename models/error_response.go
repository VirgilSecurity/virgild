package models

import (
	"fmt"
)

type ErrorResponse struct {
	Code int `json:"code"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("{\"code\":\"%v\"}", e.Code)
}

func MakeError(code int) error {
	return ErrorResponse{
		Code: code,
	}
}
