package http

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Router struct {
	card *CardController
}

func (r *Router) Handler() fasthttp.RequestHandler {
	router := fasthttprouter.New()

	router.GET("/v4/card/:id", r.card.Get)
	router.POST("/v4/card", r.card.Create)
	router.DELETE("/v4/card/:id", r.card.Revoke)
	router.POST("/card/actions/search", r.card.Search)

	return router.Handler
}
