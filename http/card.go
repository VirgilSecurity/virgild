package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgild/core"
)

type CardController struct {
	Card core.CardHandler
}

func (h *CardController) Get(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx}

	id := ctx.UserValue("id").(string)
	h.Card.Get(id, resp)
}

func (h *CardController) Create(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx}

	req := new(core.Request)
	err := json.Unmarshal(ctx.PostBody(), req)
	if err != nil {
		resp.Error(core.ErrorJSONIsInvalid)
		return
	}

	h.Card.Create(req, resp)
}

func (h *CardController) Revoke(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx}

	req := new(core.Request)
	err := json.Unmarshal(ctx.PostBody(), req)
	if err != nil {
		resp.Error(core.ErrorJSONIsInvalid)
		return
	}
	h.Card.Revoke(req, resp)
}

func (h *CardController) Search(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx}
	var c core.Criteria
	err := json.Unmarshal(ctx.PostBody(), &c)
	if err != nil {
		resp.Error(core.ErrorJSONIsInvalid)
		return
	}

	h.Card.Search(c, resp)
}
