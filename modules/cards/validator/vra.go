package validator

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
)

func ValidateVRASign(id string, pub virgilcrypto.PublicKey, crypto virgilcrypto.Crypto) func(req *virgil.SignableRequest) (bool, error) {
	return func(req *virgil.SignableRequest) (bool, error) {
		sign, ok := req.Meta.Signatures[id]
		if !ok {
			return false, core.ErrorMissVRASign
		}
		ok, err := crypto.Verify(req.Snapshot, sign, pub)
		if err != nil {
			return false, err
		}
		return ok, nil
	}
}
