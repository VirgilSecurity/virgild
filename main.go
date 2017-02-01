package main

import (
	"fmt"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/admin"
	"github.com/VirgilSecurity/virgild/modules/cards"
	"github.com/VirgilSecurity/virgild/modules/statistics"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

func main() {
	conf := config.Init("virgild.conf")

	fmt.Println("VirgilD CardID:", conf.Site.VirgilD.CardID)
	fmt.Println("VirgilD PubKey:", conf.Site.VirgilD.PublicKey)

	c := cards.Init(conf)
	s := statistics.Init(conf)
	a := admin.Init(conf)

	r := fasthttprouter.New()

	// Cards
	r.GET("/v4/card/:id", s.Middleware(c.GetCard))
	r.POST("/v4/card", s.Middleware(c.CreateCard))
	r.POST("/v1/card", s.Middleware(c.CreateCard))
	r.POST("/v4/card/actions/search", s.Middleware(c.SearchCards))
	r.DELETE("/v4/card/:id", s.Middleware(c.RevokeCard))
	r.DELETE("/v1/card/:id", s.Middleware(c.RevokeCard))
	r.GET("/api/cards/count", c.CountCards)

	// Statistics
	r.GET("/api/statistics", s.GetStatistic)
	r.GET("/api/statistics/last", s.LastActions)

	// Admin
	r.ServeFiles("/public/*filepath", "./public")
	r.GET("/api/card", a.CardInfo)
	r.GET("/api/config", a.GetConfig)
	r.POST("/api/config", a.UpdateConfig)
	r.GET("/", a.Index)

	panic(fasthttp.ListenAndServe(":8080", r.Handler))
}
