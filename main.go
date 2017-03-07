package main

import (
	"fmt"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/admin"
	"github.com/VirgilSecurity/virgild/modules/auth"
	"github.com/VirgilSecurity/virgild/modules/cards"
	"github.com/VirgilSecurity/virgild/modules/health"
	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"
)

func main() {
	conf := config.Init()

	if conf.Cards.Mode != config.CardModeCache {
		fmt.Println("VirgilD CardID:", conf.Site.VirgilD.CardID)
		fmt.Println("VirgilD PubKey:", conf.Site.VirgilD.PublicKey)
	}
	h := health.Init(conf)
	c := cards.Init(conf)
	a := admin.Init(conf)
	au := auth.Init(conf)

	r := fasthttprouter.New()

	// Cards
	r.GET("/v4/card/:id", au.Middleware(auth.PermissionGetCard, c.GetCard))
	r.POST("/v4/card", au.Middleware(auth.PermissionCreateCard, c.CreateCard))
	r.POST("/v1/card", au.Middleware(auth.PermissionCreateCard, c.CreateCard))
	r.POST("/v4/card/actions/search", au.Middleware(auth.PermissionSearchCards, c.SearchCards))
	r.DELETE("/v4/card/:id", au.Middleware(auth.PermissionRevokeCard, c.RevokeCard))
	r.DELETE("/v1/card/:id", au.Middleware(auth.PermissionRevokeCard, c.RevokeCard))

	// Admin
	r.GET("/api/card", a.CardInfo)
	r.GET("/api/tokens", a.Auth(au.GetTokens))
	r.POST("/api/tokens", a.Auth(au.CreateToken))
	r.DELETE("/api/tokens/:id", a.Auth(au.RemoveToken))
	r.PUT("/api/tokens/:id", a.Auth(au.UpdateToken))

	r.GET("/health/status", h.Status)
	r.GET("/health/info", a.Auth(h.Info))

	panic(fasthttp.ListenAndServe(conf.Common.Address, r.Handler))
}
