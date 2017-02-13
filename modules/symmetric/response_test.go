package symmetric

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type fakeLogger struct {
	mock.Mock
}

func (f *fakeLogger) Printf(format string, args ...interface{}) {
	f.Called()
}

type fakeData struct {
	Name string
}

func TestMakeResponseWrapper_ReturnSeccess(t *testing.T) {
	expected := fakeData{"Test"}
	s := MakeResponseWrapper(nil)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		return expected, nil
	})
	ctx := &fasthttp.RequestCtx{}
	s(ctx)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())

	var actual fakeData
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestMakeResponseWrapper_ResponseErrorCode_ReturnErr(t *testing.T) {
	for k, v := range code2Resp {
		expected := responseErrorModel{k, v}
		s := MakeResponseWrapper(nil)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
			return nil, k
		})
		ctx := &fasthttp.RequestCtx{}
		s(ctx)
		var actual responseErrorModel
		err := json.Unmarshal(ctx.Response.Body(), &actual)
		assert.Nil(t, err)

		assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()), fmt.Sprintf("Content type for: %v", k))
		assert.Equal(t, expected, actual, fmt.Sprintf("Content for: %v", k))
		assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode(), fmt.Sprintf("StatusCode type for: %v", k))
	}
}

func TestMakeResponseWrapper_ResponseErrorCodeUnknowCode_ReturnErr(t *testing.T) {
	expected := responseErrorModel{-1, "Unknow response error"}
	s := MakeResponseWrapper(nil)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		return nil, ResponseErrorCode(-1)
	})
	ctx := &fasthttp.RequestCtx{}
	s(ctx)
	var actual responseErrorModel
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, expected, actual)
	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode())
}

func TestMakeResponseWrapper_ResponseErrorCodeErrorEntityNotFound_ReturnErr(t *testing.T) {
	s := MakeResponseWrapper(nil)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		return nil, ErrorEntityNotFound
	})
	ctx := &fasthttp.RequestCtx{}
	s(ctx)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, fasthttp.StatusNotFound, ctx.Response.StatusCode())
}

func TestMakeResponseWrapper_NativeErr_ReturnErr(t *testing.T) {
	expected := responseErrorModel{10000, `Internal application error. You know, shit happens, so do internal server errors. Just take a deep breath and try harder.`}
	l := new(fakeLogger)
	l.On("Printf")
	s := MakeResponseWrapper(l)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		return nil, fmt.Errorf("Error")
	})
	ctx := &fasthttp.RequestCtx{}
	s(ctx)
	var actual responseErrorModel
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, expected, actual)
	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
}

func TestMakeResponseWrapper_NativeErr_LogErr(t *testing.T) {
	l := new(fakeLogger)
	l.On("Printf").Once()
	s := MakeResponseWrapper(l)(func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		return nil, fmt.Errorf("Error")
	})
	ctx := &fasthttp.RequestCtx{}
	s(ctx)

	l.AssertExpectations(t)
}
