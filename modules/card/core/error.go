package core

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
)

var (
	UnsupportedAuthTypeErr = coreapi.APIError{
		Code:       20300,
		StatusCode: http.StatusUnauthorized,
	}
	JSONInvalidErr = coreapi.APIError{
		Code:       30000,
		StatusCode: http.StatusBadRequest,
	}
	SnapshotIncorrectErr = coreapi.APIError{
		Code:       30107,
		StatusCode: http.StatusBadRequest,
	}
	CardIdentityEmptyErr = coreapi.APIError{
		Code:       30114,
		StatusCode: http.StatusBadRequest,
	}
	PublicKeyLentghInvalidErr = coreapi.APIError{
		Code:       30117,
		StatusCode: http.StatusBadRequest,
	}
	SignsIsEmptyErr = coreapi.APIError{
		Code:       30123,
		StatusCode: http.StatusBadRequest,
	}
	SingItemInvalidForApplicationErr = coreapi.APIError{
		Code:       30128,
		StatusCode: http.StatusBadRequest,
	}
	SignItemInvalidForClientErr = coreapi.APIError{
		Code:       30142,
		StatusCode: http.StatusBadRequest,
	}
	CardDataCannotContainsMoreThan16EntriesErr = coreapi.APIError{
		Code:       30103,
		StatusCode: http.StatusBadRequest,
	}
	DataValueExceed256Err = coreapi.APIError{
		Code:       30121,
		StatusCode: http.StatusBadRequest,
	}
	InfoValueExceed256Err = coreapi.APIError{
		Code:       30105,
		StatusCode: http.StatusBadRequest,
	}
	GlobalCardIdentityTypeMustBeEmailErr = coreapi.APIError{
		Code:       30100,
		StatusCode: http.StatusBadRequest,
	}
	RevocationReasonIsEmptyErr = coreapi.APIError{
		Code:       30139,
		StatusCode: http.StatusBadRequest,
	}
	RevokeCardIDInURLNotEqualCardIDInBodyErr = coreapi.APIError{
		Code:       30131,
		StatusCode: http.StatusBadRequest,
	}
	ScopeMustBeGlobalOrApplicationErr = coreapi.APIError{
		Code:       30101,
		StatusCode: http.StatusBadRequest,
	}
	SearchIdentitesEmptyErr = coreapi.APIError{
		Code:       30111,
		StatusCode: http.StatusBadRequest,
	}
	EmailIdentityIvalidErr = coreapi.APIError{
		Code:       30115,
		StatusCode: http.StatusBadRequest,
	}
	VRASignInvalidErr = coreapi.APIError{
		Code:       30143,
		StatusCode: http.StatusBadRequest,
	}
)
