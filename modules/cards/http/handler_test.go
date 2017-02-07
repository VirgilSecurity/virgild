package http

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func makeRequestCtx(body interface{}) *fasthttp.RequestCtx {
	res := &fasthttp.RequestCtx{
		Request: fasthttp.Request{
			Header: fasthttp.RequestHeader{},
		},
		Response: fasthttp.Response{},
	}

	switch body.(type) {
	case []byte:
		res.Request.AppendBody(body.([]byte))
	case string:
		res.Request.AppendBodyString(body.(string))
	case nil:

	default:
		b, _ := json.Marshal(body)
		res.Request.AppendBody(b)

	}

	return res
}

func TestGetCard_RunFunc(t *testing.T) {
	var executed bool
	const id = "test"
	f := GetCard(func(actual string) (*core.Card, error) {
		executed = true
		assert.Equal(t, id, actual)
		return new(core.Card), fmt.Errorf("Error")
	})
	ctx := fasthttp.RequestCtx{}
	ctx.SetUserValue("id", id)
	c, err := f(&ctx)

	assert.NotNil(t, err)
	assert.NotNil(t, c)
	assert.True(t, executed)
}

func TestSearchCards_JSONInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx("bo',dy")
	_, err := SearchCards(nil)(ctx)
	assert.Equal(t, core.ErrorJSONIsInvalid, err)
}

func TestSearchCards_FuncExecuted(t *testing.T) {
	crit := virgil.SearchCriteriaByAppBundle("test")
	ctx := makeRequestCtx(crit)
	var executed bool
	cs, err := SearchCards(func(c *virgil.Criteria) ([]core.Card, error) {
		assert.Equal(t, crit, c)
		executed = true
		return make([]core.Card, 0), fmt.Errorf("Error")
	})(ctx)

	assert.NotNil(t, cs)
	assert.NotNil(t, err)
	assert.True(t, executed)
}

func TestCreateCard_JSONInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx("bo',dy")
	_, err := CreateCard(nil)(ctx)
	assert.Equal(t, core.ErrorJSONIsInvalid, err)
}

func TestCreateCard_SnapshotInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx(virgil.SignableRequest{Snapshot: []byte("test")})
	_, err := CreateCard(nil)(ctx)
	assert.Equal(t, core.ErrorSnapshotIncorrect, err)
}

func TestCreateCard_FuncExecuted(t *testing.T) {
	kp, _ := virgil.Crypto().GenerateKeypair()
	epub, _ := kp.PublicKey().Encode()
	req, _ := virgil.NewCreateCardRequest("test", "nick", kp.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"test": "test",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone",
			DeviceName: "my",
		},
	})
	expected := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "nick",
			PublicKey:    epub,
			Scope:        virgil.CardScope.Application,
			Data: map[string]string{
				"test": "test",
			},
			DeviceInfo: virgil.DeviceInfo{
				Device:     "iphone",
				DeviceName: "my",
			},
		},
		Request: *req,
	}
	ctx := makeRequestCtx(req)
	var executed bool
	cs, err := CreateCard(func(r *core.CreateCardRequest) (*core.Card, error) {
		assert.Equal(t, expected, r)
		executed = true
		return &core.Card{}, fmt.Errorf("Error")
	})(ctx)

	assert.NotNil(t, cs)
	assert.NotNil(t, err)
	assert.True(t, executed)
}

func TestRevokeCard_JSONInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx("bo',dy")
	_, err := RevokeCard(nil)(ctx)
	assert.Equal(t, core.ErrorJSONIsInvalid, err)
}

func TestRevokeCard_SnapshotInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx(virgil.SignableRequest{Snapshot: []byte("test")})
	_, err := RevokeCard(nil)(ctx)
	assert.Equal(t, core.ErrorSnapshotIncorrect, err)
}

func TestRevokeCard_FuncExecuted(t *testing.T) {
	req, _ := virgil.NewRevokeCardRequest("123", virgil.RevocationReason.Unspecified)
	expected := &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			ID:               "123",
			RevocationReason: virgil.RevocationReason.Unspecified,
		},
		Request: *req,
	}
	ctx := makeRequestCtx(req)
	var executed bool
	_, err := RevokeCard(func(r *core.RevokeCardRequest) error {
		assert.Equal(t, expected, r)
		executed = true
		return fmt.Errorf("Error")
	})(ctx)

	assert.NotNil(t, err)
	assert.True(t, executed)
}

type fakeCardsRepository struct {
	mock.Mock
}

func (f *fakeCardsRepository) Count() (int64, error) {
	args := f.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestGetCountCards_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardsRepository)
	repo.On("Count").Return(int64(0), fmt.Errorf("Error"))
	s := GetCountCards(repo)
	_, err := s(&fasthttp.RequestCtx{})

	assert.NotNil(t, err)
}

func TestGetCountCards_RepoReturnCount_ReturnCount(t *testing.T) {
	repo := new(fakeCardsRepository)
	repo.On("Count").Return(int64(10), nil)
	s := GetCountCards(repo)
	cm, err := s(&fasthttp.RequestCtx{})

	assert.Nil(t, err)
	assert.Equal(t, cardsCountModel{10}, cm)
}
