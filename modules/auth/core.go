package auth

import "fmt"

type errResponse int

func (e errResponse) Error() string {
	return fmt.Sprintf("code: %v", e)
}

var (
	errTokenInvalid     errResponse = 20300
	errAuthServiceDenny errResponse = 20302
	errForbidden        errResponse = 20500
)
