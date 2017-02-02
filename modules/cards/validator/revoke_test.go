package validator

import (
	"fmt"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
)

func TestRevokeCard_ReasonInvalide_ReturnErr(t *testing.T) {
	r := RevokeCard(func(req *core.RevokeCardRequest) error {
		return nil
	})
	req := new(core.RevokeCardRequest)

	err := r(req)

	assert.Equal(t, core.ErrorRevocationReasonIsEmpty, err)
}

func TestRevokeCard_AddCustomValidator_Executed(t *testing.T) {
	var executed bool
	r := RevokeCard(func(req *core.RevokeCardRequest) error {
		return nil
	}, func(req *core.RevokeCardRequest) (bool, error) {
		executed = true
		return true, nil
	})
	req := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			RevocationReason: virgil.RevocationReason.Compromised,
			ID:               "1234",
		},
	}

	r(req)

	assert.True(t, executed)
}

func TestRevokeCard_ReqMissId_ReturnErr(t *testing.T) {
	r := RevokeCard(func(req *core.RevokeCardRequest) error {
		return nil
	})
	req := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			RevocationReason: virgil.RevocationReason.Compromised,
		},
	}

	err := r(req)

	assert.Equal(t, core.ErrorRevokeCardIDInURLNotEqualCardIDInBody, err)
}

func TestRevokeCard_ReqValid_NextFuncExecuted(t *testing.T) {
	var executed bool
	r := RevokeCard(func(req *core.RevokeCardRequest) error {
		executed = true
		return nil
	})
	req := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			RevocationReason: virgil.RevocationReason.Compromised,
			ID:               "1234",
		},
	}

	r(req)

	assert.True(t, executed)
}

func TestWrapRevokeValidateVRASign_Executed(t *testing.T) {
	var executed bool
	w := WrapRevokeValidateVRASign(func(req *virgil.SignableRequest) (bool, error) {
		executed = true
		return true, fmt.Errorf("Error")
	})
	ok, err := w(new(core.RevokeCardRequest))

	assert.True(t, executed)
	assert.True(t, ok)
	assert.NotNil(t, err)
}
