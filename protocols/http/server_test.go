package http

import (
	"fmt"
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
		return d, protocols.CodeResponse(args.Get(1).(int))
	} else {
		return nil, protocols.CodeResponse(args.Get(1).(int))
	}
}

func (c *MockController) SearchCards(data []byte) ([]byte, protocols.CodeResponse) {
	args := c.Called(data)
	if d, ok := args.Get(0).([]byte); ok {
		return d, protocols.CodeResponse(args.Get(1).(int))
	} else {
		return nil, protocols.CodeResponse(args.Get(1).(int))
	}
}

func (c *MockController) CreateCard(data []byte) ([]byte, protocols.CodeResponse) {
	args := c.Called(data)
	if d, ok := args.Get(0).([]byte); ok {
		return d, protocols.CodeResponse(args.Get(1).(int))
	} else {
		return nil, protocols.CodeResponse(args.Get(1).(int))
	}
}

func (c *MockController) RevokeCard(id string, data []byte) protocols.CodeResponse {
	args := c.Called(id, data)
	return protocols.CodeResponse(args.Get(1).(int))
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
	h.SetMethod(method)
	h.SetRequestURI("http://test.com" + path)
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

	testTable := map[string]string{
		"/v4/card":                "POST",
		"/v4/card/actions/search": "POST",
		"/v4/card/asdf":           "GET",
		"/v4/card/test":           "DELETE",
	}
	for k, v := range testTable {
		ctx := MakeContext(v, k, "test")
		r.router.HandleRequest(ctx)
		assert.Equal(t, 401, ctx.Response.StatusCode())
		assert.Equal(t, expected, ctx.Response.Body())
	}
}

func Test_Auth_TokenCorrect_ReturnErr(t *testing.T) {
	id := "test"
	expected := []byte(`expected result`)

	r, c, h := MakeRouter()
	h.On("Auth", "test").Return(true, nil)
	c.On("GetCard", id).Return(expected, protocols.Ok)
	ctx := MakeContext("GET", "/v4/card/test", "test")

	r.router.HandleRequest(ctx)

	assert.Equal(t, 200, ctx.Response.StatusCode())
	assert.Equal(t, expected, ctx.Response.Body())
}

func Test_GetCard_ControllerReturnNotFound_BodyEmpty(t *testing.T) {
	id := "test"
	r, c, h := MakeRouter()
	h.On("Auth", "test").Return(true, nil)
	c.On("GetCard", id).Return(nil, protocols.NotFound)
	ctx := MakeContext("GET", "/v4/card/test", "test")

	r.router.HandleRequest(ctx)

	assert.Equal(t, 404, ctx.Response.StatusCode())
	assert.Nil(t, ctx.Response.Body())
}

func Test_GetCard_ControllerCheckMappingStatusCode(t *testing.T) {
	id := "test"
	testTable := map[protocols.CodeResponse]int{
		protocols.RequestError: 400,
		// protocols.Ok:           200,
		// protocols.ServerError:  500,
	}
	for k, v := range testTable {
		fmt.Println("Key:", k, "Val:", v)
		r, c, h := MakeRouter()
		h.On("Auth", "test").Return(true, nil)
		c.On("GetCard", id).Return(nil, k)
		ctx := MakeContext("GET", "/v4/card/test", "test")
		r.router.HandleRequest(ctx)

		assert.Equal(t, v, ctx.Response.StatusCode())
	}

}
