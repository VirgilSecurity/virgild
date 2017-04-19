package coreapi

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"-"`
	StatusCode int    `json:"-"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("code:%v msg:%s", err.Code, err.Message)
}

var InternalServerErr = APIError{
	Code:       10000,
	StatusCode: http.StatusInternalServerError,
}
