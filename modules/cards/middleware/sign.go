package middleware

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
)

func SignCreateRequest(signer func(req *virgil.SignableRequest) error, next core.CreateCard) core.CreateCard {
	return func(req *core.CreateCardRequest) (*core.Card, error) {
		err := signer(&req.Request)
		if err != nil {
			return nil, err
		}
		return next(req)
	}
}

func SignRevokeRequest(signer func(req *virgil.SignableRequest) error, next core.RevokeCard) core.RevokeCard {
	return func(req *core.RevokeCardRequest) error {
		err := signer(&req.Request)
		if err != nil {
			return err
		}
		return next(req)
	}
}

func MakeSigner(cardID string, priv virgilcrypto.PrivateKey) func(req *virgil.SignableRequest) error {
	signer := virgil.RequestSigner{}
	return func(req *virgil.SignableRequest) error {
		err := signer.AuthoritySign(req, cardID, priv)
		if err != nil {
			return err
		}
		return nil
	}
}
