package symmetric

import "fmt"

type ResponseErrorCode int

func (code ResponseErrorCode) Error() string {
	return fmt.Sprintf("code: %v", int(code))
}

const (
	// 401
	ErrorAuthHeaderInvalid ResponseErrorCode = 20300

	// 403
	ErrorForbidden ResponseErrorCode = 20500

	// 404
	ErrorEntityNotFound ResponseErrorCode = 404

	// 400
	ErrorJSONIsInvalid ResponseErrorCode = 30000
)
