package storage

import "fmt"

var (
	ErrorNotFound  = fmt.Errorf("Card not found")
	ErrorForbidden = fmt.Errorf("Request forbidden")
)
