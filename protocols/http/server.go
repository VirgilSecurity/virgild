package http

import (
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
)

func MakeServer(host string, c protocols.Controller, auth protocols.AuthHandler) protocols.Server {
	r := router{
		controller:  c,
		authHandler: auth,
	}
	r.init()
	return &server{
		router: r,
		host:   host,
	}
}

type server struct {
	host   string
	router router
}

func (s *server) Serve() error {
	return fasthttp.ListenAndServe(s.host, s.router.router.HandleRequest)
}

type router struct {
	router      *routing.Router
	controller  protocols.Controller
	authHandler protocols.AuthHandler
}

func (r *router) init() {
	r.router = routing.New()
	v4 := r.router.Group("/v4")
	v4.Use(r.auth)

	v4.Get("/card/<id>", r.getCard)
	v4.Post("/card", r.createCard)
	v4.Post("/card/actions/search", r.search)
	v4.Delete("/card/<id>", r.delete)

}
func (r *router) auth(ctx *routing.Context) error {
	b := ctx.Request.Header.Peek("Authorization")
	isAuth, data := r.authHandler.Auth(string(b[:]))
	fmt.Println("We are here")
	if isAuth {
		fmt.Println("Auth")
		ctx.Next()
	} else {
		fmt.Println("Not Auth")
		ctx.SetStatusCode(401)
		ctx.Write(data)
	}
	return nil
}

func (r *router) getCard(ctx *routing.Context) error {
	id := ctx.Param("id")
	res, code := r.controller.GetCard(id)
	setStatus(ctx, code)
	if code != protocols.NotFound {
		ctx.Write(res)
	}
	return nil
}
func (r *router) createCard(ctx *routing.Context) error {
	data := ctx.PostBody()
	res, code := r.controller.CreateCard(data)

	setStatus(ctx, code)
	ctx.Write(res)
	return nil
}

func (r *router) search(ctx *routing.Context) error {
	data := ctx.PostBody()
	res, code := r.controller.SearchCards(data)
	setStatus(ctx, code)
	ctx.Write(res)
	return nil
}

func (r *router) delete(ctx *routing.Context) error {
	data := ctx.PostBody()
	id := ctx.Param("id")
	code := r.controller.RevokeCard(id, data)
	setStatus(ctx, code)
	return nil
}

func setStatus(ctx *routing.Context, code protocols.CodeResponse) {
	switch code {
	case protocols.Ok:
		ctx.SetStatusCode(200)
	case protocols.NotFound:
		ctx.SetStatusCode(404)
	case protocols.RequestError:
		ctx.SetStatusCode(400)
	case protocols.ServerError:
		ctx.SetStatusCode(500)
	}
}
