package middleware

import (
	"context"
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func makeFakeSignerfunc(err error) signerFunc {
	return func(req *virgil.SignableRequest) error {
		return err
	}
}

func TestMakeSigner_SignerReturnErr_ReturnErr(t *testing.T) {
	req := &virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}}
	kp, _ := virgil.Crypto().GenerateKeypair()
	signerMiddleware := MakeSigner("123", kp.PrivateKey())
	err := signerMiddleware(req)
	assert.Nil(t, err)

	ok, _ := virgil.Crypto().Verify(virgil.Crypto().CalculateFingerprint(req.Snapshot), req.Meta.Signatures["123"], kp.PublicKey())

	assert.True(t, ok)
}

func TestSignCreateCardRequest_SignerReturnErr_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	signer := makeFakeSignerfunc(fmt.Errorf("Error"))
	var executed bool
	signerMiddleware := SignCreateRequest(signer, func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		executed = true
		return nil, nil
	})

	_, err := signerMiddleware(context.Background(), req)
	assert.NotNil(t, err)
	assert.False(t, executed)
}

func TestSignCreateCardRequest_NextFuncExecuted(t *testing.T) {
	req := &core.CreateCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}

	signer := makeFakeSignerfunc(nil)
	var executed bool
	signerMiddleware := SignCreateRequest(signer, func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		executed = true
		return nil, nil
	})

	signerMiddleware(context.Background(), req)
	assert.True(t, executed)
}

func TestSignRevokeRequest_SignerReturnErr_ReturnErr(t *testing.T) {
	req := &core.RevokeCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	signer := makeFakeSignerfunc(fmt.Errorf("Error"))
	var executed bool
	signerMiddleware := SignRevokeRequest(signer, func(ctx context.Context, req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})

	err := signerMiddleware(context.Background(), req)
	assert.NotNil(t, err)
	assert.False(t, executed)
}

func TestSignRevokeRequest_NextFuncExecuted(t *testing.T) {
	req := &core.RevokeCardRequest{
		Request: virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}},
	}
	signer := makeFakeSignerfunc(nil)
	var executed bool
	signerMiddleware := SignRevokeRequest(signer, func(ctx context.Context, req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})

	signerMiddleware(context.Background(), req)
	assert.True(t, executed)
}

func TestMakeSigner_UnsuppoertedKey_ReturnErr(t *testing.T) {
	signer := MakeSigner("id", nil)
	err := signer(&virgil.SignableRequest{Snapshot: []byte("test"), Meta: virgil.RequestMeta{Signatures: make(map[string][]byte)}})
	assert.Error(t, err)
}
