package http

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Logger interface {
	Println(v ...interface{})
	Panicf(format string, v ...interface{})
}

type Controller interface {
	GetCard(id string) ([]byte, error)
}

func MakeRouter(contreoller Controller, logger Logger) Router {
	router := Router{
		controller: contreoller,
		logger:     logger,
		router:     routing.New(),
	}
	router.Init()
	return router
}

type Router struct {
	controller Controller
	logger     Logger
	router     *routing.Router
}

func (r *Router) GetHandleRequest() fasthttp.RequestHandler {
	return r.router.HandleRequest
}

func (r *Router) Init() {
	v4 := r.router.Group("/v4")

	v4.Get("/card/<id>", func(ctx *routing.Context) error {
		id := ctx.Param("id")
		res, err := r.controller.GetCard(id)

		if err != nil {
			return err
			r.logger.Println("Get card by id (", id, ") Error:", err)
			ctx.Error("Service unavailable. Please contact your administrator.", 500)
		}

		ctx.Write(res)
		return nil
	})
}
