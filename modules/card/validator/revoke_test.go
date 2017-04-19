package validator

import (
	"context"
	"fmt"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
)

func makeContext(id string) context.Context {
	return core.SetURLCardID(context.Background(), id)
}

func makeFakeRevokeCardRequest() *core.RevokeCardRequest {
	req, _ := virgil.NewRevokeCardRequest("id", virgil.RevocationReason.Compromised)
	req.AppendSignature("id", []byte("signature"))
	return &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			ID:               "id",
			RevocationReason: virgil.RevocationReason.Compromised,
		},
		Request: *req,
	}
}

func TestRevokeCard_ReasonInvalide_ReturnErr(t *testing.T) {
	r := RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	})
	req := new(core.RevokeCardRequest)

	err := r(context.Background(), req)

	assert.Equal(t, core.RevocationReasonIsEmptyErr, err)
}

func TestRevokeCard_AddCustomValidator_Executed(t *testing.T) {
	var executed bool
	r := RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	}, func(ctx context.Context, req *core.RevokeCardRequest) (bool, error) {
		executed = true
		return true, nil
	})
	req := makeFakeRevokeCardRequest()
	r(makeContext(req.Info.ID), req)
	assert.True(t, executed)
}

func TestRevokeCard_ReqMissId_ReturnErr(t *testing.T) {
	r := RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	})
	req := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			RevocationReason: virgil.RevocationReason.Compromised,
			ID:               "1234",
		},
	}

	err := r(context.Background(), req)

	assert.Equal(t, core.RevokeCardIDInURLNotEqualCardIDInBodyErr, err)
}

func TestRevokeCard_ReqValid_NextFuncExecuted(t *testing.T) {
	var executed bool
	r := RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})
	req := makeFakeRevokeCardRequest()

	r(makeContext(req.Info.ID), req)

	assert.True(t, executed)
}

func TestRevokeCard_SignEmpty_ReturnErr(t *testing.T) {
	r := RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	})
	req := makeFakeRevokeCardRequest()
	req.Request.Meta.Signatures = make(map[string][]byte)

	err := r(makeContext(req.Info.ID), req)

	assert.Equal(t, core.SignsIsEmptyErr, err)
}

func TestWrapRevokeValidateVRASign_Executed(t *testing.T) {
	var executed bool
	wrap := WrapRevokeValidateVRASign(func(req *virgil.SignableRequest) (bool, error) {
		executed = true
		return true, fmt.Errorf("Error")
	})
	ok, err := wrap(context.Background(), new(core.RevokeCardRequest))

	assert.True(t, executed)
	assert.True(t, ok)
	assert.NotNil(t, err)
}
