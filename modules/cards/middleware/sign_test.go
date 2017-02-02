package middleware

import (
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestMakeSigner_SignerReturnErr_ReturnErr(t *testing.T) {
	req := &virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}}
	kp, _ := virgil.Crypto().GenerateKeypair()
	s := MakeSigner("123", kp.PrivateKey())
	err := s(req)
	assert.Nil(t, err)

	ok, _ := virgil.Crypto().Verify(virgil.Crypto().CalculateFingerprint(req.Snapshot), req.Meta.Signatures["123"], kp.PublicKey())

	assert.True(t, ok)
}

func TestSignCreateCardRequest_SignerReturnErr_ReturnErr(t *testing.T) {
	req := core.CreateCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	var executed bool
	s := SignCreateRequest(func(req *virgil.SignableRequest) error {
		return fmt.Errorf("Error")
	}, func(req *core.CreateCardRequest) (*core.Card, error) {
		executed = true
		return nil, nil
	})

	_, err := s(&req)
	assert.NotNil(t, err)
	assert.False(t, executed)
}

func TestSignCreateCardRequest_NextFuncExecuted(t *testing.T) {
	req := core.CreateCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	var executed bool
	s := SignCreateRequest(func(req *virgil.SignableRequest) error {
		return nil
	}, func(req *core.CreateCardRequest) (*core.Card, error) {
		executed = true
		return nil, nil
	})

	s(&req)
	assert.True(t, executed)
}

func TestSignRevokeRequest_SignerReturnErr_ReturnErr(t *testing.T) {
	req := core.RevokeCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	var executed bool
	s := SignRevokeRequest(func(req *virgil.SignableRequest) error {
		return fmt.Errorf("Error")
	}, func(req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})

	err := s(&req)
	assert.NotNil(t, err)
	assert.False(t, executed)
}

func TestSignRevokeRequest_NextFuncExecuted(t *testing.T) {
	req := core.RevokeCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	var executed bool
	s := SignRevokeRequest(func(req *virgil.SignableRequest) error {
		return nil
	}, func(req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})

	s(&req)
	assert.True(t, executed)
}
