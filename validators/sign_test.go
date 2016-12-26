package validators

import (
	"github.com/stretchr/testify/assert"
	"github.com/virgilsecurity/virgild/models"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"testing"
)

var errExpected = models.MakeError(30137)

func Test_Validate_SignMiss_ReturnErr(t *testing.T) {
	r := &models.CardResponse{
		Snapshot: []byte(`Test`),
	}
	crypto := virgil.Crypto()
	keyPair, _ := crypto.GenerateKeypair()
	pub, _ := keyPair.PublicKey().Encode()
	v := MakeSignValidator(map[string][]byte{
		"id": pub,
	})
	err := v.Validate(r)
	assert.Equal(t, errExpected, err)
}

func Test_Validate_SignCorrect_ReturnNil(t *testing.T) {
	r := &models.CardResponse{
		Snapshot: []byte(`Test`),
		Meta: models.ResponseMeta{
			Signatures: make(map[string][]byte, 0),
		},
	}
	crypto := virgil.Crypto()
	keyPair, _ := crypto.GenerateKeypair()
	fp := crypto.CalculateFingerprint(r.Snapshot)
	sign, _ := crypto.Sign(fp, keyPair.PrivateKey())
	r.Meta.Signatures["id"] = sign

	pub, _ := keyPair.PublicKey().Encode()
	v := MakeSignValidator(map[string][]byte{
		"id": pub,
	})
	err := v.Validate(r)
	assert.Nil(t, err)
}

func Test_Validate_SignIncorrect_ReturnErr(t *testing.T) {
	r := &models.CardResponse{
		Snapshot: []byte(`Test`),
		Meta: models.ResponseMeta{
			Signatures: make(map[string][]byte, 0),
		},
	}
	crypto := virgil.Crypto()
	keyPair, _ := crypto.GenerateKeypair()
	r.Meta.Signatures["id"] = []byte(`test`)

	pub, _ := keyPair.PublicKey().Encode()
	v := MakeSignValidator(map[string][]byte{
		"id": pub,
	})
	err := v.Validate(r)
	assert.Equal(t, errExpected, err)
}
