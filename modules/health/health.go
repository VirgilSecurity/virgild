package health

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type info func() (map[string]interface{}, error)

type healthResult struct {
	Name   string
	Status int
	Info   map[string]interface{}
}

type HealthChecker struct {
	checkList map[string]info
}

func (h *HealthChecker) Status(ctx *fasthttp.RequestCtx) {
	r := h.check()
	for _, v := range r {
		if v.Status != fasthttp.StatusOK {
			ctx.SetStatusCode(v.Status)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *HealthChecker) Info(ctx *fasthttp.RequestCtx) {
	resp := make(map[string]interface{})
	r := h.check()
	ctx.SetStatusCode(fasthttp.StatusOK)
	for _, v := range r {
		if v.Status != fasthttp.StatusOK {
			ctx.SetStatusCode(v.Status)
		}
		info := v.Info
		if info == nil {
			info = make(map[string]interface{})
		}
		info["status"] = v.Status
		resp[v.Name] = info
	}
	b, _ := json.Marshal(resp)
	ctx.Write(b)
}

func (h *HealthChecker) check() []healthResult {
	r := make([]healthResult, 0)
	for k, v := range h.checkList {
		m := healthResult{}
		m.Name = k
		m.Status = fasthttp.StatusOK

		info, err := v()
		if err != nil {
			m.Status = fasthttp.StatusBadRequest
		}
		m.Info = info
		r = append(r, m)
	}
	return r
}
