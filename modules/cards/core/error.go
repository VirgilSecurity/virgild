package core

import "fmt"

type ResponseErrorCode int

func (code ResponseErrorCode) Error() string {
	return fmt.Sprintf("code: %v", code)
}

const (
	// 500
	ErrorInernalApplication ResponseErrorCode = 10000

	// 401
	ErrorAuthHeaderInvalid ResponseErrorCode = 20300

	// 403
	ErrorForbidden ResponseErrorCode = 20500

	// 404
	ErrorEntityNotFound ResponseErrorCode = 404

	// 400
	ErrorJSONIsInvalid                           ResponseErrorCode = 30000
	ErrorGlobalCardIdentityTypeMustBeEmail       ResponseErrorCode = 30100
	ErrorScopeMustBeGlobalOrApplication          ResponseErrorCode = 30101
	ErrorCardIDValidationFailed                  ResponseErrorCode = 30102 // miss
	ErrorCardDataCannotContainsMoreThan16Entries ResponseErrorCode = 30103
	ErrorEmptyCardInfo                           ResponseErrorCode = 30104 // miss
	ErrorInfoValueExceed256                      ResponseErrorCode = 30105
	ErrorCardDataMustBeAssotiativeArray          ResponseErrorCode = 30106 // ErrorSnapshotIncorrect
	ErrorSnapshotIncorrect                       ResponseErrorCode = 30107
	ErrorSearchIdentitesEmpty                    ResponseErrorCode = 30111
	ErrorCardIdentityTypeInvalid                 ResponseErrorCode = 30113 // miss
	ErrorCardIdentityEmpty                       ResponseErrorCode = 30114
	ErrorCardIdentityEmailInvalid                ResponseErrorCode = 30115 // miss
	ErrorCardIdentityApplicationInvalid          ResponseErrorCode = 30116 // miss
	ErrorPublicKeyLentghInvalid                  ResponseErrorCode = 30117
	ErrorPublicKeyMustBeBase64                   ResponseErrorCode = 30118 // ErrorSnapshotIncorrect
	ErrorDataParamMustBeKeyValue                 ResponseErrorCode = 30119 // ErrorCardDataMustBeAssotiativeArray, ErrorDataOrInfoValueExceed256
	ErrorDataParamMustBeString                   ResponseErrorCode = 30120 // ErrorCardDataMustBeAssotiativeArray, ErrorDataOrInfoValueExceed256
	ErrorDataValueExceed256                      ResponseErrorCode = 30121
	ErrorSignsIsEmpty                            ResponseErrorCode = 30123
	ErrorApplicationSignIsInvalid                ResponseErrorCode = 30128 // miss
	ErrorRevokeCardIDInURLNotEqualCardIDInBody   ResponseErrorCode = 30131 // miss
	ErrorCardDataParamIsAplphanumerical          ResponseErrorCode = 30134 // miss
	ErrorMissSignIdentityService                 ResponseErrorCode = 30137 // miss
	ErrorCardAlreadyExists                       ResponseErrorCode = 30138 // miss
	ErrorRevocationReasonIsEmpty                 ResponseErrorCode = 30139
	ErrorSignInvalid                             ResponseErrorCode = 30140 // miss
	ErrorSignerIsNotFound                        ResponseErrorCode = 30141 // miss
	ErrorSignItemIsInvalid                       ResponseErrorCode = 30142 // ErrorSignInvalid
	ErrorMissVRASign                             ResponseErrorCode = 30143 // miss
)
