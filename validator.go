package main

import (
	"strings"

	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type RequestValidator struct {
	SearchValidators     []func(criteria *virgil.Criteria) (bool, error)
	CreateCardValidators []func(req *CreateCardRequest) (bool, error)
	RevokeCardValidators []func(req *RevokeCardRequest) (bool, error)
}

func (v *RequestValidator) IsValidSearchCriteria(criteria *virgil.Criteria) (bool, error) {
	for _, v := range v.SearchValidators {
		if ok, err := v(criteria); !ok {
			return false, err
		}
	}
	return true, nil
}

func (v *RequestValidator) IsValidCreateCardRequest(req *CreateCardRequest) (bool, error) {
	for _, v := range v.CreateCardValidators {
		if ok, err := v(req); !ok {
			return false, err
		}
	}
	return true, nil
}

func (v *RequestValidator) IsValidRevokeCardRequest(req *RevokeCardRequest) (bool, error) {
	for _, v := range v.RevokeCardValidators {
		if ok, err := v(req); !ok {
			return false, err
		}
	}
	return true, nil
}

func ScopeMustGlobalOrApplication(criteria *virgil.Criteria) (bool, error) {
	if criteria.Scope == virgil.CardScope.Application || criteria.Scope == virgil.CardScope.Global {
		return true, nil
	}
	return false, ErrorScopeMustBeGlobalOrApplication
}

func SearchIdentitiesNotEmpty(crit *virgil.Criteria) (bool, error) {
	if len(crit.Identities) == 0 {
		return false, ErrorSearchIdentitesEmpty
	}
	return true, nil
}

func CardIdentityIsEmpty(req *CreateCardRequest) (bool, error) {
	if len(req.Info.Identity) == 0 {
		return false, ErrorCardIdentityEmpty
	}
	return true, nil
}

func CardPublicKeyLengthInvalid(req *CreateCardRequest) (bool, error) {
	if len(req.Info.PublicKey) < 16 || len(req.Info.PublicKey) > 2048 {
		return false, ErrorPublicKeyLentghInvalid
	}
	return true, nil
}

func CreateCardRequestSignsEmpty(req *CreateCardRequest) (bool, error) {
	if len(req.Request.Signatures) == 0 {
		return false, ErrorSignsIsEmpty
	}
	return true, nil
}

func CardDataEntries(req *CreateCardRequest) (bool, error) {
	if len(req.Info.Data) > 16 {
		return false, ErrorCardDataCannotContainsMoreThan16Entries
	}
	return true, nil
}

func CardDataValueExceed256(req *CreateCardRequest) (bool, error) {
	for _, v := range req.Info.Data {
		if len(v) > 256 {
			return false, ErrorDataValueExceed256
		}
	}
	return true, nil
}

func CardInfoValueExceed256(req *CreateCardRequest) (bool, error) {
	if len(req.Info.DeviceInfo.Device) > 256 || len(req.Info.DeviceInfo.DeviceName) > 256 {
		return false, ErrorInfoValueExceed256
	}
	return true, nil
}

func GlobalCardIdentityTypeMustBeEmail(req *CreateCardRequest) (bool, error) {
	if req.Info.Scope == virgil.CardScope.Global {
		if strings.ToLower(req.Info.IdentityType) != "email" {
			return false, ErrorGlobalCardIdentityTypeMustBeEmail
		}
	}
	return true, nil
}

func RevokeCardRequestSignsEmpty(req *RevokeCardRequest) (bool, error) {
	if len(req.Request.Signatures) == 0 {
		return false, ErrorSignsIsEmpty
	}
	return true, nil
}

func RevocationReasonIsInvalide(req *RevokeCardRequest) (bool, error) {
	if len(req.Info.RevocationReason) == 0 {
		return false, ErrorRevocationReasonIsEmpty
	}
	return true, nil
}
