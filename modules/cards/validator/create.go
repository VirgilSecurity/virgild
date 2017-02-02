package validator

import (
	"encoding/hex"
	"strings"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

var createValidator = []func(req *core.CreateCardRequest) (bool, error){
	cardIdentityIsEmpty,
	globalCardIdentityTypeMustBeEmail,
	cardPublicKeyLengthInvalid,
	cardDataEntries,
	cardDataValueExceed256,
	cardInfoValueExceed256,
	createCardRequestSignsEmpty,
	createCardRequestSelfSignIvalid,
}

func CreateCard(next core.CreateCard, validators ...func(req *core.CreateCardRequest) (bool, error)) core.CreateCard {
	validators = append(createValidator, validators...)
	return func(req *core.CreateCardRequest) (*core.Card, error) {
		for _, v := range validators {
			if ok, err := v(req); !ok {
				return nil, err
			}
		}
		return next(req)
	}
}

func cardIdentityIsEmpty(req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.Identity) == 0 {
		return false, core.ErrorCardIdentityEmpty
	}
	return true, nil
}

func cardPublicKeyLengthInvalid(req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.PublicKey) < 16 || len(req.Info.PublicKey) > 2048 {
		return false, core.ErrorPublicKeyLentghInvalid
	}
	return true, nil
}

func createCardRequestSignsEmpty(req *core.CreateCardRequest) (bool, error) {
	if len(req.Request.Meta.Signatures) == 0 {
		return false, core.ErrorSignsIsEmpty
	}
	return true, nil
}

func createCardRequestSelfSignIvalid(req *core.CreateCardRequest) (bool, error) {
	id := virgil.Crypto().CalculateFingerprint(req.Request.Snapshot)
	sign, ok := req.Request.Meta.Signatures[hex.EncodeToString(id)]
	if !ok {
		return false, core.ErrorSignItemInvalidForClient
	}
	pub, err := virgil.Crypto().ImportPublicKey(req.Info.PublicKey)
	if err != nil {
		return false, core.ErrorSnapshotIncorrect
	}
	if ok, err = virgil.Crypto().Verify(id, sign, pub); !ok {
		return false, core.ErrorSignItemInvalidForClient
	}
	return true, nil
}

func cardDataEntries(req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.Data) > 16 {
		return false, core.ErrorCardDataCannotContainsMoreThan16Entries
	}
	return true, nil
}

func cardDataValueExceed256(req *core.CreateCardRequest) (bool, error) {
	for _, v := range req.Info.Data {
		if len(v) > 256 {
			return false, core.ErrorDataValueExceed256
		}
	}
	return true, nil
}

func cardInfoValueExceed256(req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.DeviceInfo.Device) > 256 || len(req.Info.DeviceInfo.DeviceName) > 256 {
		return false, core.ErrorInfoValueExceed256
	}
	return true, nil
}

func globalCardIdentityTypeMustBeEmail(req *core.CreateCardRequest) (bool, error) {
	if req.Info.Scope == virgil.CardScope.Global {
		if strings.ToLower(req.Info.IdentityType) != "email" {
			return false, core.ErrorGlobalCardIdentityTypeMustBeEmail
		}
	}
	return true, nil
}

func WrapCreateValidateVRASign(f func(req *virgil.SignableRequest) (bool, error)) func(req *core.CreateCardRequest) (bool, error) {
	return func(req *core.CreateCardRequest) (bool, error) {
		return f(&req.Request)
	}
}
