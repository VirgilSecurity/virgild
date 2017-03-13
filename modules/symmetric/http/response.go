package http

import (
	"encoding/json"

	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func MakeResponseWrapper(logger Logger) func(f core.Response) fasthttp.RequestHandler {
	return func(f core.Response) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			s, err := f(ctx)

			ctx.SetContentType("application/json")
			ctx.ResetBody()

			if err != nil {
				responseError(err, ctx, logger)
			} else {
				responseSeccess(s, ctx, logger)
			}
		}
	}
}

type responseErrorModel struct {
	Code    core.ResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

func responseError(err error, ctx *fasthttp.RequestCtx, logger Logger) {
	initErr := errors.Cause(err)
	switch e := initErr.(type) {
	case core.ResponseErrorCode:
		if e == core.ErrorEntityNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusBadRequest)

		json.NewEncoder(ctx).Encode(responseErrorModel{
			Code:    e,
			Message: mapCode2Msg(e),
		})
	default:
		logger.Printf("Intrenal error: %+v", err)

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		json.NewEncoder(ctx).Encode(responseErrorModel{
			Code:    core.ResponseErrorCode(10000),
			Message: mapCode2Msg(10000),
		})
	}
}

func responseSeccess(model interface{}, ctx *fasthttp.RequestCtx, logger Logger) {
	err := json.NewEncoder(ctx).Encode(model)
	if err != nil {
		responseError(err, ctx, logger)
	}
}

var code2Resp = map[core.ResponseErrorCode]string{
	10000: `Internal application error. You know, shit happens, so do internal server errors. Just take a deep breath and try harder.`,
	30000: `JSON specified as a request is invalid`,
}

func mapCode2Msg(code core.ResponseErrorCode) string {
	if msg, ok := code2Resp[code]; ok {
		return msg
	}
	return "Unknow response error"
}
