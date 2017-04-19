package validator

import (
	"context"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	virgil "gopkg.in/virgil.v4"
)

type revokeCardValidatorHandler func(ctx context.Context, req *core.RevokeCardRequest) (bool, error)

var revokeValidator = []revokeCardValidatorHandler{
	revocationReasonIsInvalide,
	revokeCardMissId,
	revokeCardRequestSignsEmpty,
}

func RevokeCard(f core.RevokeCardHandler, validators ...revokeCardValidatorHandler) core.RevokeCardHandler {
	validators = append(revokeValidator, validators...)
	return func(ctx context.Context, req *core.RevokeCardRequest) error {
		for _, v := range validators {
			if ok, err := v(ctx, req); !ok {
				return err
			}
		}
		return f(ctx, req)
	}
}

func revokeCardRequestSignsEmpty(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
	if len(req.Request.Meta.Signatures) == 0 {
		return false, core.SignsIsEmptyErr
	}
	return true, nil
}

func revocationReasonIsInvalide(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
	if len(req.Info.RevocationReason) == 0 {
		return false, core.RevocationReasonIsEmptyErr
	}
	return true, nil
}

func revokeCardMissId(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
	urlID := core.GetURLCardID(ctx)
	if req.Info.ID != urlID {
		return false, core.RevokeCardIDInURLNotEqualCardIDInBodyErr
	}
	return true, nil
}

func WrapRevokeValidateVRASign(f func(req *virgil.SignableRequest) (bool, error)) func(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
	return func(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
		return f(&req.Request)
	}
}
