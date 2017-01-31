package validator

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

var revokeValidator = []func(req *core.RevokeCardRequest) (bool, error){
	revocationReasonIsInvalide,
	//revokeCardRequestSignsEmpty,
}

func RevokeCard(next core.RevokeCard, validators ...func(req *core.RevokeCardRequest) (bool, error)) core.RevokeCard {
	validators = append(revokeValidator, validators...)
	return func(req *core.RevokeCardRequest) error {
		for _, v := range validators {
			if ok, err := v(req); !ok {
				return err
			}
		}
		return next(req)
	}
}

func revokeCardRequestSignsEmpty(req *core.RevokeCardRequest) (bool, error) {
	if len(req.Request.Meta.Signatures) == 0 {
		return false, core.ErrorSignsIsEmpty
	}
	return true, nil
}

func revocationReasonIsInvalide(req *core.RevokeCardRequest) (bool, error) {
	if len(req.Info.RevocationReason) == 0 {
		return false, core.ErrorRevocationReasonIsEmpty
	}
	return true, nil
}

func WrapRevokeValidateVRASign(f func(req *virgil.SignableRequest) (bool, error)) func(req *core.RevokeCardRequest) (bool, error) {
	return func(req *core.RevokeCardRequest) (bool, error) {
		return f(&req.Request)
	}
}
