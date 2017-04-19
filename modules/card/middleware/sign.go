package middleware

import (
	"context"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
)

type signerFunc func(req *virgil.SignableRequest) error

func SignCreateRequest(signer signerFunc, f core.CreateCardHandler) core.CreateCardHandler {
	return func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		err := signer(&req.Request)
		if err != nil {
			return nil, err
		}
		return f(ctx, req)
	}
}

func SignRevokeRequest(signer signerFunc, f core.RevokeCardHandler) core.RevokeCardHandler {
	return func(ctx context.Context, req *core.RevokeCardRequest) error {
		err := signer(&req.Request)
		if err != nil {
			return err
		}
		return f(ctx, req)
	}
}

func MakeSigner(cardID string, priv virgilcrypto.PrivateKey) signerFunc {
	signer := virgil.RequestSigner{}
	return func(req *virgil.SignableRequest) error {
		err := signer.AuthoritySign(req, cardID, priv)
		if err != nil {
			return err
		}
		return nil
	}
}
