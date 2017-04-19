package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	virgil "gopkg.in/virgil.v4"
)

type brokenReader struct{}

func (r brokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("ERROR")
}

func TestGetCard_SetContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.RawQuery = ":id=test_id"

	h := func(ctx context.Context, id string) (*virgil.CardResponse, error) {
		assert.Equal(t, "test_id", id)

		return nil, nil
	}

	GetCard(h)(req)
}

func TestSearchCards_BodyBroken_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", brokenReader{})

	h := func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := SearchCards(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestSearchCards_BodyInvalidJSON_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := SearchCards(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestSearchCards_JSONCorrectPars(t *testing.T) {
	expected := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		IdentityType: "test",
		Scope:        virgil.CardScope.Application,
	}
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"identities":["bob","alice"],"identity_type":"test","scope":"application"}`))

	h := func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		assert.Equal(t, expected, crit)

		return nil, nil
	}

	SearchCards(h)(req)
}

func TestCreateCard_BodyBroken_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", brokenReader{})

	h := func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := CreateCard(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestCreateCard_BodyInvalidJSON_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := CreateCard(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestCreateCard_SnapshotInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJwdWJsaWMiOiJNQ293QlFZREsyVndBeUVBTTh6UGx6VUNWOC9SVGREQXVDeXBmenR0V280OFZ0U0k2YUVCYTdPcnEwYz0iLCJpZGVudGl0eSI6MTQ5MjQ5NjE4MDczNkBtYWlsaW5hdG9yLmNvbSwiaWRlbnRpdHlfdHlwZSI6ImVtYWlsIiwic2NvcGUiOiJnbG9iYWwiLCJpbmZvIjp7ImRldmljZSI6Im1hYyIsImRldmljZV9uYW1lIjoibWFjYm9vayBwcm8ifSwiZGF0YSI6eyJ2aXJnaWxfYXV0b3Rlc3QiOiJ2aXJnaWxfYXV0b3Rlc3QifX0=",
  "meta": {"signs": { } }}`))

	h := func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := CreateCard(h)(req)

	assert.Equal(t, core.SnapshotIncorrectErr, err)
}

func TestCreateCard_JSONCorrectPars(t *testing.T) {
	expected := &core.CreateCardRequest{
		Info: virgil.CardModel{
			PublicKey:    []byte(`public_key`),
			Identity:     "is email",
			IdentityType: "email",
			Scope:        virgil.CardScope.Global,
			DeviceInfo: virgil.DeviceInfo{
				Device:     "mac",
				DeviceName: "macbook pro",
			},
			Data: map[string]string{
				"virgil_data": "virgil_data",
			},
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte(`{"public_key":"cHVibGljX2tleQ==","identity":"is email","identity_type":"email","scope":"global","info":{"device":"mac","device_name":"macbook pro"},"data":{"virgil_data":"virgil_data"}}`),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": []byte(`:)`),
				},
				Validation: &virgil.ValidationInfo{
					Token: "validation token",
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJwdWJsaWNfa2V5IjoiY0hWaWJHbGpYMnRsZVE9PSIsImlkZW50aXR5IjoiaXMgZW1haWwiLCJpZGVudGl0eV90eXBlIjoiZW1haWwiLCJzY29wZSI6Imdsb2JhbCIsImluZm8iOnsiZGV2aWNlIjoibWFjIiwiZGV2aWNlX25hbWUiOiJtYWNib29rIHBybyJ9LCJkYXRhIjp7InZpcmdpbF9kYXRhIjoidmlyZ2lsX2RhdGEifX0=",
  "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" },"validation": { "token": "validation token" } }}`))

	h := func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		assert.Equal(t, expected, req)

		return nil, nil
	}

	CreateCard(h)(req)
}

func TestRevokeCard_BodyBroken_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", brokenReader{})

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	}

	_, err := RevokeCard(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestRevokeCard_BodyInvalidJSON_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	}

	_, err := RevokeCard(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestRevokeCard_SnapshotInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJwdWJsaWMiOiJNQ293QlFZREsyVndBeUVBTTh6UGx6VUNWOC9SVGREQXVDeXBmenR0V280OFZ0U0k2YUVCYTdPcnEwYz0iLCJpZGVudGl0eSI6MTQ5MjQ5NjE4MDczNkBtYWlsaW5hdG9yLmNvbSwiaWRlbnRpdHlfdHlwZSI6ImVtYWlsIiwic2NvcGUiOiJnbG9iYWwiLCJpbmZvIjp7ImRldmljZSI6Im1hYyIsImRldmljZV9uYW1lIjoibWFjYm9vayBwcm8ifSwiZGF0YSI6eyJ2aXJnaWxfYXV0b3Rlc3QiOiJ2aXJnaWxfYXV0b3Rlc3QifX0=",
  "meta": {"signs": { } }}`))

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	}

	_, err := RevokeCard(h)(req)

	assert.Equal(t, core.SnapshotIncorrectErr, err)
}

func TestRevokeCard_JSONCorrectPars(t *testing.T) {
	expected := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			ID:               "1234",
			RevocationReason: virgil.RevocationReason.Compromised,
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte(`{"card_id":"1234","revocation_reason":"compromised"}`),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": []byte(`:)`),
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJjYXJkX2lkIjoiMTIzNCIsInJldm9jYXRpb25fcmVhc29uIjoiY29tcHJvbWlzZWQifQ==",
  "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" } }}`))

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		assert.Equal(t, expected, req)

		return nil
	}

	RevokeCard(h)(req)
}

func TestRevokeCard_ReturnEmptyObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJjYXJkX2lkIjoiMTIzNCIsInJldm9jYXRpb25fcmVhc29uIjoiY29tcHJvbWlzZWQifQ==",
  "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" } }}`))

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	}

	seccess, _ := RevokeCard(h)(req)

	assert.Equal(t, []byte(`{}`), seccess)
}

func TestRevokeCard_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJjYXJkX2lkIjoiMTIzNCIsInJldm9jYXRpb25fcmVhc29uIjoiY29tcHJvbWlzZWQifQ==",
  "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" } }}`))

	h := func(ctx context.Context, req *core.RevokeCardRequest) error {
		return fmt.Errorf("ERROR")
	}

	_, err := RevokeCard(h)(req)

	assert.Error(t, err, "ERROR")
}

func TestCreateRelation_BodyBroken_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", brokenReader{})

	h := func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := CreateRelation(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestCreateRelation_BodyInvalidJSON_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := CreateRelation(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestCreateRelation_JSONCorrectPars(t *testing.T) {
	expected := &core.CreateRelationRequest{
		ID: "1234",
		Request: virgil.SignableRequest{
			Snapshot: []byte(`{"public_key":"cHVibGljX2tleQ==","identity":"is email","identity_type":"email","scope":"global","info":{"device":"mac","device_name":"macbook pro"},"data":{"virgil_data":"virgil_data"}}`),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": []byte(`:)`),
				},
				Validation: &virgil.ValidationInfo{
					Token: "validation token",
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJwdWJsaWNfa2V5IjoiY0hWaWJHbGpYMnRsZVE9PSIsImlkZW50aXR5IjoiaXMgZW1haWwiLCJpZGVudGl0eV90eXBlIjoiZW1haWwiLCJzY29wZSI6Imdsb2JhbCIsImluZm8iOnsiZGV2aWNlIjoibWFjIiwiZGV2aWNlX25hbWUiOiJtYWNib29rIHBybyJ9LCJkYXRhIjp7InZpcmdpbF9kYXRhIjoidmlyZ2lsX2RhdGEifX0=",
  "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" },"validation": { "token": "validation token" } }}`))

	req.URL.RawQuery = ":id=1234"
	h := func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		assert.Equal(t, expected, req)

		return nil, nil
	}

	CreateRelation(h)(req)
}

func TestRevokeRelation_BodyBroken_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", brokenReader{})

	h := func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := RevokeRelation(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestRevokeRelation_BodyInvalidJSON_ReturnErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := RevokeRelation(h)(req)

	assert.Equal(t, core.JSONInvalidErr, err)
}

func TestRevokeRelation_SnapshotInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
  "content_snapshot": "eyJwdWJsaWMiOiJNQ293QlFZREsyVndBeUVBTTh6UGx6VUNWOC9SVGREQXVDeXBmenR0V280OFZ0U0k2YUVCYTdPcnEwYz0iLCJpZGVudGl0eSI6MTQ5MjQ5NjE4MDczNkBtYWlsaW5hdG9yLmNvbSwiaWRlbnRpdHlfdHlwZSI6ImVtYWlsIiwic2NvcGUiOiJnbG9iYWwiLCJpbmZvIjp7ImRldmljZSI6Im1hYyIsImRldmljZV9uYW1lIjoibWFjYm9vayBwcm8ifSwiZGF0YSI6eyJ2aXJnaWxfYXV0b3Rlc3QiOiJ2aXJnaWxfYXV0b3Rlc3QifX0=",
  "meta": {"signs": { } }}`))

	h := func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return nil, nil
	}

	_, err := RevokeRelation(h)(req)

	assert.Equal(t, core.SnapshotIncorrectErr, err)
}

func TestRevokeRelation_JSONCorrectPars(t *testing.T) {
	expected := &core.RevokeRelationRequest{
		ID: "4321",
		Info: virgil.RevokeCardRequest{
			ID:               "1234",
			RevocationReason: virgil.RevocationReason.Compromised,
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte(`{"card_id":"1234","revocation_reason":"compromised"}`),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": []byte(`:)`),
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
    "content_snapshot": "eyJjYXJkX2lkIjoiMTIzNCIsInJldm9jYXRpb25fcmVhc29uIjoiY29tcHJvbWlzZWQifQ==",
    "meta": { "signs": {"c84ba35ed7af45948495659ce0bbdfd0db82f745d9e1436548474edd1bcf7a75": "Oik=" } }}`))
	req.URL.RawQuery = ":id=4321"

	h := func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		assert.Equal(t, expected, req)

		return nil, nil
	}

	RevokeRelation(h)(req)
}
