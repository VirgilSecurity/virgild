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
	SearchCards([]byte) ([]byte, error)
	CreateCard([]byte) ([]byte, error)
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

	v4.Post("/card", func(ctx *routing.Context) error {
		data := ctx.PostBody()
		res, err := r.controller.CreateCard(data)

		if err != nil {
			return err
		}

		ctx.Write(res)
		return nil
	})

	v4.Post("/card/actions/search", func(ctx *routing.Context) error {
		data := ctx.PostBody()
		res, err := r.controller.SearchCards(data)

		if err != nil {
			return err
		}

		ctx.Write(res)
		return nil
	})

	v4.Get("/card/<id>", func(ctx *routing.Context) error {
		id := ctx.Param("id")
		res, err := r.controller.GetCard(id)

		if err != nil {
			return err
		}

		ctx.Write(res)
		return nil
	})

}
