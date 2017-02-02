package validator

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
)

func ValidateVRASign(id string, pub virgilcrypto.PublicKey) func(req *virgil.SignableRequest) (bool, error) {
	crypto := virgil.Crypto()
	return func(req *virgil.SignableRequest) (bool, error) {
		sign, ok := req.Meta.Signatures[id]
		if !ok {
			return false, core.ErrorVRASignInvalid
		}
		ok, err := crypto.Verify(virgil.Crypto().CalculateFingerprint(req.Snapshot), sign, pub)
		if err != nil {
			return false, core.ErrorVRASignInvalid
		}
		return ok, nil
	}
}
