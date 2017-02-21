package main

import (
	"fmt"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/admin"
	"github.com/VirgilSecurity/virgild/modules/auth"
	"github.com/VirgilSecurity/virgild/modules/cards"
	"github.com/VirgilSecurity/virgild/modules/statistics"
	"github.com/VirgilSecurity/virgild/modules/symmetric"
	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"
)

func main() {
	conf := config.Init()

	if conf.Cards.Mode != config.CardModeCache {
		fmt.Println("VirgilD CardID:", conf.Site.VirgilD.CardID)
		fmt.Println("VirgilD PubKey:", conf.Site.VirgilD.PublicKey)
	}

	c := cards.Init(conf)
	s := statistics.Init(conf)
	a := admin.Init(conf)
	au := auth.Init(conf)
	sk := symmetric.Init(conf)

	r := fasthttprouter.New()

	// Cards
	r.GET("/v4/card/:id", au.Middleware(auth.PermissionGetCard, s.Middleware(c.GetCard)))
	r.POST("/v4/card", au.Middleware(auth.PermissionCreateCard, s.Middleware(c.CreateCard)))
	r.POST("/v1/card", au.Middleware(auth.PermissionCreateCard, s.Middleware(c.CreateCard)))
	r.POST("/v4/card/actions/search", au.Middleware(auth.PermissionSearchCards, s.Middleware(c.SearchCards)))
	r.DELETE("/v4/card/:id", au.Middleware(auth.PermissionRevokeCard, s.Middleware(c.RevokeCard)))
	r.DELETE("/v1/card/:id", au.Middleware(auth.PermissionRevokeCard, s.Middleware(c.RevokeCard)))

	// symmetric keys
	r.POST("/api/keys", sk.CreateKey)
	r.GET("/api/users/:user_id", sk.GetKeysForUser)
	r.GET("/api/keys/:key_id/users", sk.GetUsersForKey)
	r.GET("/api/keys/:key_id/users/:user_id", sk.GetKey)

	// Admin
	r.ServeFiles("/public/*filepath", "./public")
	r.GET("/api/cards/count", a.Auth(c.CountCards))
	r.GET("/api/card", a.Auth(a.CardInfo))
	r.GET("/api/config", a.Auth(a.GetConfig))
	r.POST("/api/config", a.Auth(a.UpdateConfig))
	r.GET("/api/tokens", a.Auth(au.GetTokens))
	r.POST("/api/tokens", a.Auth(au.CreateToken))
	r.DELETE("/api/tokens/:id", a.Auth(au.RemoveToken))
	r.PUT("/api/tokens/:id", a.Auth(au.UpdateToken))
	r.GET("/", a.Index)

	// Statistics
	r.GET("/api/statistics", a.Auth(s.GetStatistic))
	r.GET("/api/statistics/last", a.Auth(s.LastActions))

	panic(fasthttp.ListenAndServe(conf.Common.Address, r.Handler))
}
