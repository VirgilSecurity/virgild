package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgild/core"
)

type response struct {
	ctx *fasthttp.RequestCtx
}

func (r *response) Error(err core.ResponseError) {
	r.ctx.SetContentType("application/json")

	r.ctx.ResetBody()
	status := fasthttp.StatusBadRequest
	if err == core.ErrorInernalApplication {
		status = fasthttp.StatusInternalServerError
	}
	r.ctx.SetStatusCode(status)
	json.NewEncoder(r.ctx).Encode(err)
}

func (r *response) Success(model interface{}) {
	r.ctx.SetContentType("application/json")
	err := json.NewEncoder(r.ctx).Encode(model)
	if err != nil {
		r.Error(core.ErrorInernalApplication)
		return
	}
}
