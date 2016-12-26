package signer

import (
	"github.com/virgilsecurity/virgild/models"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/virgilcrypto"
)

type ServiceSigner struct {
	ID         string
	PrivateKey virgilcrypto.PrivateKey
}

func (s *ServiceSigner) Sign(r *models.CardResponse) error {
	crypto := virgil.Crypto()
	sign, err := crypto.Sign(r.Snapshot, s.PrivateKey)
	if r.Meta.Signatures == nil {
		r.Meta.Signatures = make(map[string][]byte, 0)
	}
	r.Meta.Signatures[s.ID] = sign
	return err
	// return nil
}
