package validator

import (
	"testing"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestValidateVRASign_SignNotExist_ReturnErr(t *testing.T) {
	req := &virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}}
	kp, _ := virgil.Crypto().GenerateKeypair()
	s := ValidateVRASign("123", kp.PublicKey())

	_, err := s(req)
	assert.Equal(t, core.ErrorVRASignInvalid, err)
}

func TestValidateVRASign_SignInvalid_ReturnErr(t *testing.T) {
	req := &virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}}
	kp, _ := virgil.Crypto().GenerateKeypair()
	req.Meta.Signatures["123"] = []byte("123")
	s := ValidateVRASign("123", kp.PublicKey())

	_, err := s(req)
	assert.Equal(t, core.ErrorVRASignInvalid, err)
}

func TestValidateVRASign_SignValid_ReturnErr(t *testing.T) {
	req := &virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}}
	kp, _ := virgil.Crypto().GenerateKeypair()
	signer := virgil.RequestSigner{}
	signer.AuthoritySign(req, "123", kp.PrivateKey())
	s := ValidateVRASign("123", kp.PublicKey())

	ok, _ := s(req)
	assert.True(t, ok)
}
