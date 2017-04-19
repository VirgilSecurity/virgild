package validator

import (
	"context"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/VirgilSecurity/virgild/modules/card/core"

	virgil "gopkg.in/virgil.v4"
)

type createValidatorHandler func(ctx context.Context, req *core.CreateCardRequest) (bool, error)

var createValidator = []createValidatorHandler{
	cardIdentityIsEmpty,
	globalCardIdentityTypeMustBeEmail,
	emailCardIdentityInvalid,
	cardPublicKeyLengthInvalid,
	cardDataEntries,
	cardDataValueExceed256,
	cardInfoValueExceed256,
	createCardRequestSignsEmpty,
	createCardRequestSelfSignIvalid,
}

func CreateCard(f core.CreateCardHandler, validators ...createValidatorHandler) core.CreateCardHandler {
	validators = append(createValidator, validators...)
	return func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		for _, v := range validators {
			if ok, err := v(ctx, req); !ok {
				return nil, err
			}
		}
		return f(ctx, req)
	}
}

func cardIdentityIsEmpty(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.Identity) == 0 {
		return false, core.CardIdentityEmptyErr
	}
	return true, nil
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func emailCardIdentityInvalid(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if req.Info.IdentityType == "email" && !emailRegexp.MatchString(req.Info.Identity) {
		return false, core.EmailIdentityIvalidErr
	}
	return true, nil
}

func cardPublicKeyLengthInvalid(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.PublicKey) < 16 || len(req.Info.PublicKey) > 2048 {
		return false, core.PublicKeyLentghInvalidErr
	}
	return true, nil
}

func createCardRequestSignsEmpty(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if len(req.Request.Meta.Signatures) == 0 {
		return false, core.SignsIsEmptyErr
	}
	return true, nil
}

func createCardRequestSelfSignIvalid(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	id := virgil.Crypto().CalculateFingerprint(req.Request.Snapshot)
	sign, ok := req.Request.Meta.Signatures[hex.EncodeToString(id)]
	if !ok {
		return false, core.SignItemInvalidForClientErr
	}
	pub, err := virgil.Crypto().ImportPublicKey(req.Info.PublicKey)
	if err != nil {
		return false, core.SnapshotIncorrectErr
	}
	if ok, err = virgil.Crypto().Verify(id, sign, pub); !ok {
		return false, core.SignItemInvalidForClientErr
	}
	return true, nil
}

func cardDataEntries(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.Data) > 16 {
		return false, core.CardDataCannotContainsMoreThan16EntriesErr
	}
	return true, nil
}

func cardDataValueExceed256(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	for _, v := range req.Info.Data {
		if len(v) > 256 {
			return false, core.DataValueExceed256Err
		}
	}
	return true, nil
}

func cardInfoValueExceed256(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if len(req.Info.DeviceInfo.Device) > 256 || len(req.Info.DeviceInfo.DeviceName) > 256 {
		return false, core.InfoValueExceed256Err
	}
	return true, nil
}

func globalCardIdentityTypeMustBeEmail(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	if req.Info.Scope == virgil.CardScope.Global {
		if strings.ToLower(req.Info.IdentityType) != "email" {
			return false, core.GlobalCardIdentityTypeMustBeEmailErr
		}
	}
	return true, nil
}

func WrapCreateValidateVRASign(f func(req *virgil.SignableRequest) (bool, error)) func(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
	return func(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
		return f(&req.Request)
	}
}
