package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
)

func TestRequestOwner_AuthHeaderEmpty_ContextIsNotSet(t *testing.T) {
	var token string
	var authHeader string

	h := func(req *http.Request) (interface{}, error) {
		token = core.GetOwnerRequest(req.Context())
		authHeader = core.GetAuthHeader(req.Context())

		return nil, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	RequestOwner(h)(req)

	assert.Equal(t, "", token)
	assert.Equal(t, "", authHeader)
}

func TestRequestOwner_AuthHeaderInvalid_ReturnErr(t *testing.T) {

	h := func(req *http.Request) (interface{}, error) {
		return nil, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "berrer 1234")
	_, err := RequestOwner(h)(req)

	assert.Equal(t, core.UnsupportedAuthTypeErr, err)
}

func TestRequestOwner_AuthHeaderExist_ContextIsSet(t *testing.T) {
	var token string
	var authHeader string

	h := func(req *http.Request) (interface{}, error) {
		token = core.GetOwnerRequest(req.Context())
		authHeader = core.GetAuthHeader(req.Context())

		return nil, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "VIRGIL 1234")
	RequestOwner(h)(req)

	assert.Equal(t, "1234", token)
	assert.Equal(t, "VIRGIL 1234", authHeader)
}
