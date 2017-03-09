package health

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func makeInfoFunc(r map[string]interface{}, err error) info {
	return func() (map[string]interface{}, error) {
		return r, err
	}
}

func TestStatus_CheckerInfo_ReturnOk(t *testing.T) {
	h := HealthChecker{}
	ctx := makeRequestCtx(nil)
	h.Status(ctx)

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
}

func TestInfo_CheckerReturnErr_ReturnBadRequest(t *testing.T) {
	expected, _ := json.Marshal(map[string]interface{}{
		"test": map[string]interface{}{
			"status": fasthttp.StatusBadRequest,
		},
	})
	h := HealthChecker{checkList: map[string]info{
		"test": makeInfoFunc(nil, fmt.Errorf("error")),
	}}
	ctx := makeRequestCtx(nil)
	h.Info(ctx)

	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	assert.Equal(t, expected, ctx.Response.Body())
}

func TestInfo_CheckerInfo_ReturnOk(t *testing.T) {
	expected, _ := json.Marshal(map[string]interface{}{
		"test": map[string]interface{}{
			"status":  fasthttp.StatusOK,
			"latency": 12,
		},
	})
	h := HealthChecker{checkList: map[string]info{
		"test": makeInfoFunc(map[string]interface{}{"latency": 12}, nil),
	}}
	ctx := makeRequestCtx(nil)
	h.Info(ctx)

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.Equal(t, expected, ctx.Response.Body())
}
