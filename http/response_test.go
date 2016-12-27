package http

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgild/core"
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

func assertResponse(t *testing.T, expected core.ResponseError, r *fasthttp.RequestCtx) {
	statusCpde := fasthttp.StatusBadRequest
	if expected == core.ErrorInernalApplication {
		statusCpde = fasthttp.StatusInternalServerError
	}
	//  else if expected == core.StatusErrorAttemptNotFound {
	// 	statusCpde = fasthttp.StatusNotFound
	// }
	assert.Equal(t, statusCpde, r.Response.StatusCode())

	// if expected == core.StatusErrorAttemptNotFound {
	// 	return
	// }

	var err core.ResponseError
	json.Unmarshal(r.Response.Body(), &err)
	assert.Equal(t, expected, err)
}

func TestError_StatusErrorInernalApplicationError_Return500(t *testing.T) {
	ctx := makeRequestCtx("body")
	resp := response{ctx: ctx}
	resp.Error(core.ErrorInernalApplication)
	assertResponse(t, core.ErrorInernalApplication, ctx)
}

// func TestError_StatusErrorAttemptNotFound_Return404(t *testing.T) {
// 	ctx := makeRequestCtx("body")
// 	resp := response{ctx: ctx}
// 	resp.Error(core.StatusErrorAttemptNotFound)
// 	assertResponse(t, core.StatusErrorAttemptNotFound, ctx)
// }
//
//
// func TestError_OtherErrors_Return400(t *testing.T) {
// 	table := []core.ResponseStatus{
// 		core.StatusErrorUUIDValidFailed,
// 		core.StatusErrorCardNotFound,
// 		core.StatusErrorEncryptedMessageValidationFailed,
// 		core.StatusErrorAccessTokenExpired,
// 		core.StatusErrorUnsupportedGrantType,
// 		core.StatusErrorCodeNotFound,
// 		core.StatusErrorCodeExpired,
// 		core.StatusErrorCodeWasUsed,
// 		core.StatusErrorAccessTokenBroken,
// 		core.StatusErrorRefreshTokenNotFound,
// 		core.StatusErrorCardMustGlobal,
// 	}
// 	for _, v := range table {
// 		ctx := makeRequestCtx("body")
// 		resp := response{ctx: ctx}
// 		resp.Error(v)
// 		assertResponse(t, v, ctx)
// 	}
// }

func TestSuccess(t *testing.T) {
	ctx := makeRequestCtx("body")
	resp := response{ctx: ctx}
	resp.Success("str")

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))

	// New line in the end of bytes will be added response body
	assert.Equal(t, []byte("\"str\"\n"), ctx.Response.Body())
}
