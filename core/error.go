package core

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	// 500
	ErrorInernalApplication = ResponseError{
		Code:    10000,
		Message: `Internal application error. You know, shit happens, so do internal server errors. Just take a deep breath and try harder.`,
	}

	// 403
	ErrorForbidden = ResponseError{
		Code:    20500,
		Message: `The Virgil Card is not available in this application`,
	}
	// 400
	ErrorJSONIsInvalid = ResponseError{
		Code:    30000,
		Message: `JSON specified as a request is invalid`,
	}
)
