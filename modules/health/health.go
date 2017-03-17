package health

import (
	"github.com/rcrowley/go-metrics"
	"github.com/valyala/fasthttp"
)

type info func() (map[string]interface{}, error)

type healthResult struct {
	Name   string
	Status int
	Info   map[string]interface{}
}

type HealthChecker struct {
}

func (h *HealthChecker) Status(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *HealthChecker) Info(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	metrics.WriteJSONOnce(metrics.DefaultRegistry, ctx)
}
