package validators

import (
	"github.com/virgilsecurity/virgild/models"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/virgilcrypto"
)

func MakeSignValidator(keys map[string][]byte) *SignValidator {
	validator := SignValidator{
		keys: make(map[string]virgilcrypto.PublicKey, 0),
	}
	for k, v := range keys {
		pub, err := virgilcrypto.DecodePublicKey(v)
		if err != nil {

		}
		validator.keys[k] = pub
	}
	return &validator
}

type SignValidator struct {
	keys map[string]virgilcrypto.PublicKey
}

func (v *SignValidator) Validate(r *models.CardResponse) *models.ErrorResponse {
	crypto := virgil.Crypto()
	fp := crypto.CalculateFingerprint(r.Snapshot)
	for id, key := range v.keys {
		sign, ok := r.Meta.Signatures[id]
		if !ok {
			return models.MakeError(30137)
		}

		valid, _ := crypto.Verify(fp, sign, key)
		if !valid {
			return models.MakeError(30137)
		}
	}
	return nil
}
