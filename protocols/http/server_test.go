package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
	"testing"
)

type MockController struct {
	mock.Mock
}

func (c *MockController) GetCard(id string) ([]byte, protocols.CodeResponse) {
	args := c.Called(id)
	if d, ok := args.Get(0).([]byte); ok {
		return d, args.Get(1).(protocols.CodeResponse)
	} else {
		return nil, args.Get(1).(protocols.CodeResponse)
	}
}

func (c *MockController) SearchCards(data []byte) ([]byte, protocols.CodeResponse) {
	args := c.Called(data)
	if d, ok := args.Get(0).([]byte); ok {
		return d, args.Get(1).(protocols.CodeResponse)
	} else {
		return nil, args.Get(1).(protocols.CodeResponse)
	}
}

func (c *MockController) CreateCard(data []byte) ([]byte, protocols.CodeResponse) {
	args := c.Called(data)
	if d, ok := args.Get(0).([]byte); ok {
		return d, args.Get(1).(protocols.CodeResponse)
	} else {
		return nil, args.Get(1).(protocols.CodeResponse)
	}
}

func (c *MockController) RevokeCard(id string, data []byte) protocols.CodeResponse {
	args := c.Called(id, data)
	return args.Get(1).(protocols.CodeResponse)
}

type MockAuthHandler struct {
	mock.Mock
}

func (h *MockAuthHandler) Auth(token string) (bool, []byte) {
	args := h.Called(token)
	if d, ok := args.Get(1).([]byte); ok {
		return args.Bool(0), d
	} else {
		return args.Bool(0), nil
	}
}

func MakeRouter() (*router, *MockController, *MockAuthHandler) {
	c := new(MockController)
	h := new(MockAuthHandler)
	r := &router{
		authHandler: h,
		controller:  c,
	}
	r.init()
	return r, c, h
}

func MakeContext(method, path, token string) *fasthttp.RequestCtx {
	h := fasthttp.RequestHeader{}
	h.Add("Authorization", token)
	return &fasthttp.RequestCtx{
		Request: fasthttp.Request{
			Header: h,
		},
		Response: fasthttp.Response{},
	}
}

func Test_Auth_TokenIncorrect_ReturnErr(t *testing.T) {
	r, _, h := MakeRouter()
	expected := []byte(`some error`)
	h.On("Auth", "test").Return(false, expected)
	ctx := MakeContext("Post", "/v4/card", "test")
	r.router.HandleRequest(ctx)
	assert.Equal(t, 401, ctx.Response.StatusCode())
}
