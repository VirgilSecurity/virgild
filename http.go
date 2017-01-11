package main

import (
	"encoding/json"

	"github.com/buaazp/fasthttprouter"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type response struct {
	ctx    *fasthttp.RequestCtx
	logger Logger
}

var code2Resp = map[ResponseErrorCode]string{
	10000: `Internal application error. You know, shit happens, so do internal server errors. Just take a deep breath and try harder.`,
	20500: `The Virgil Card is not available in this application`,
	30000: `JSON specified as a request is invalid`,
}

func mapCode2Msg(code ResponseErrorCode) string {
	if msg, ok := code2Resp[code]; ok {
		return msg
	}
	return "Unknow response error"
}

type responseError struct {
	Code    ResponseErrorCode `json:"code"`
	Message string            `json:"message"`
}

func (r *response) Error(err error) {
	r.ctx.SetContentType("application/json")
	r.ctx.ResetBody()

	initErr := errors.Cause(err)
	switch e := initErr.(type) {
	case virgil.SDKError:
		respStatus := e.TransportStatusCode()
		r.ctx.SetStatusCode(respStatus)
		if respStatus == fasthttp.StatusNotFound {
			return
		}
		code := ResponseErrorCode(e.ServiceStatusCode())
		json.NewEncoder(r.ctx).Encode(responseError{
			Code:    code,
			Message: mapCode2Msg(code),
		})
	case ResponseErrorCode:
		if e == ErrorEntityNotFound {
			r.ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		r.ctx.SetStatusCode(fasthttp.StatusBadRequest)
		if e == ErrorApplicationSignIsInvalid {
			r.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		json.NewEncoder(r.ctx).Encode(responseError{
			Code:    e,
			Message: mapCode2Msg(e),
		})
	default:
		r.logger.Printf("Intrenal error: %+v", err)

		r.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		json.NewEncoder(r.ctx).Encode(responseError{
			Code:    ResponseErrorCode(ErrorInernalApplication),
			Message: mapCode2Msg(ErrorInernalApplication),
		})
	}
}

func (r *response) Success(model interface{}) {
	r.ctx.SetContentType("application/json")
	err := json.NewEncoder(r.ctx).Encode(model)
	if err != nil {
		r.Error(err)
	}
}

func (r *response) Response(seccess interface{}, err error) {
	if err != nil {
		r.Error(err)
	} else {
		r.Success(seccess)
	}
}

type signableRequest struct {
	Meta struct {
		Signatures map[string][]byte `json:"signs"`
	} `json:"meta"`
	Snapshot []byte `json:"content_snapshot"`
}

type Logger interface {
	Printf(format string, args ...interface{})
}

type CardController struct {
	Card   CardHandler
	Logger Logger
}

func (h *CardController) Get(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx, logger: h.Logger}

	id := ctx.UserValue("id").(string)
	resp.Response(h.Card.Get(id))
}

func (h *CardController) Create(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx, logger: h.Logger}

	req := new(signableRequest)
	err := json.Unmarshal(ctx.PostBody(), req)
	if err != nil {
		resp.Error(ErrorJSONIsInvalid)
		return
	}
	creq := new(virgil.CreateCardRequest)
	err = json.Unmarshal(req.Snapshot, creq)
	if err != nil {
		resp.Error(ErrorSnapshotIncorrect)
		return
	}

	s, err := h.Card.Create(&CreateCardRequest{
		Info: *creq,
		Request: virgil.SignableRequest{
			Snapshot:   req.Snapshot,
			Signatures: req.Meta.Signatures,
		}})

	resp.Response(s, err)
}

func (h *CardController) Revoke(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx, logger: h.Logger}

	req := new(signableRequest)
	err := json.Unmarshal(ctx.PostBody(), req)
	if err != nil {
		resp.Error(ErrorJSONIsInvalid)
		return
	}

	creq := virgil.RevokeCardRequest{}
	err = json.Unmarshal(req.Snapshot, &creq)
	if err != nil {
		resp.Error(ErrorSnapshotIncorrect)
		return
	}

	s, err := h.Card.Revoke(&RevokeCardRequest{
		Info: creq,
		Request: virgil.SignableRequest{
			Snapshot:   req.Snapshot,
			Signatures: req.Meta.Signatures,
		}})
	resp.Response(s, err)
}

func (h *CardController) Search(ctx *fasthttp.RequestCtx) {
	resp := &response{ctx: ctx, logger: h.Logger}

	c := new(virgil.Criteria)
	err := json.Unmarshal(ctx.PostBody(), c)
	if err != nil {
		resp.Error(ErrorJSONIsInvalid)
		return
	}
	if c.Scope == "" {
		c.Scope = virgil.CardScope.Application
	}

	resp.Response(h.Card.Search(c))
}

type Router struct {
	Card *CardController
}

func (r *Router) Handler() fasthttp.RequestHandler {
	router := fasthttprouter.New()

	router.GET("/v4/card/:id", r.Card.Get)
	router.POST("/v4/card", r.Card.Create)
	router.DELETE("/v4/card/:id", r.Card.Revoke)
	router.POST("/v4/card/actions/search", r.Card.Search)

	return router.Handler
}
