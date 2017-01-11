package main

import (
	"fmt"
	"time"

	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type RequestSigner interface {
	Sign(*virgil.SignableRequest) error
}

type Fingerprint interface {
	Calculate(data []byte) string
}

type VirgilClient interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(virgil.Criteria) ([]*virgil.Card, error)
	CreateCard(req *virgil.SignableRequest) (*virgil.Card, error)
	RevokeCard(req *virgil.SignableRequest) error
}

type Validator interface {
	IsValidSearchCriteria(criteria *virgil.Criteria) (bool, error)
	IsValidCreateCardRequest(req *CreateCardRequest) (bool, error)
	IsValidRevokeCardRequest(req *RevokeCardRequest) (bool, error)
}
type CardRepository interface {
	Get(id string) (*cardSql, error)
	Find(identitis []string, identityType string, scope string) ([]cardSql, error)
	Add(cs cardSql) error
	DeleteById(id string) error
	DeleteBySearch(identitis []string, identityType string, scope string) error
}

type cardSql struct {
	CardID       string `xorm:"Index 'card_id'"`
	Identity     string `xorm:"Index"`
	IdentityType string
	Scope        string
	ExpireAt     time.Time
	Deleted      bool
	ErrorCode    int
	Card         []byte
}
type Card struct {
	ID       string   `json:"id"`
	Snapshot []byte   `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     CardMeta `json:"meta"`
}
type CardMeta struct {
	CreatedAt   string            `json:"created_at,omitempty"`
	CardVersion string            `json:"card_version,omitempty"`
	Signatures  map[string][]byte `json:"signs"`
}

type CalculateFingerprint interface {
	CalculateFingerprint(data []byte) []byte
}

type ResponseErrorCode int

func (code ResponseErrorCode) Error() string {
	return fmt.Sprintf("code: %v", code)
}

const (
	// 500
	ErrorInernalApplication ResponseErrorCode = 10000

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

type CardHandler interface {
	Get(id string) (interface{}, error)
	Search(criteria *virgil.Criteria) (interface{}, error)
	Create(req *CreateCardRequest) (interface{}, error)
	Revoke(req *RevokeCardRequest) (interface{}, error)
}

type CreateCardRequest struct {
	Info    virgil.CreateCardRequest
	Request virgil.SignableRequest
}

type RevokeCardRequest struct {
	Info    virgil.RevokeCardRequest
	Request virgil.SignableRequest
}
