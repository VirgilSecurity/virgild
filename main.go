package main

import (
	"fmt"
	"net"
	"time"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/admin"
	"github.com/VirgilSecurity/virgild/modules/auth"
	"github.com/VirgilSecurity/virgild/modules/cards"
	"github.com/VirgilSecurity/virgild/modules/health"
	"github.com/buaazp/fasthttprouter"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"

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

	f := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			t := metrics.GetOrRegisterTimer("response", nil)
			t.Time(func() {
				next(ctx)
			})
		}
	}

	if conf.Common.Config.Metrics.Log.Enabled {
		go metrics.LogScaled(metrics.DefaultRegistry, conf.Common.Config.Metrics.Log.Interval, time.Microsecond, conf.Common.Logger)
	}
	if conf.Common.Config.Metrics.Graphite.Address != "" {
		addr, _ := net.ResolveTCPAddr("tcp", conf.Common.Config.Metrics.Graphite.Address)
		graphanaConf := graphite.Config{
			Addr:          addr,
			Registry:      metrics.DefaultRegistry,
			FlushInterval: conf.Common.Config.Metrics.Graphite.Interval,
			DurationUnit:  time.Microsecond,
			Prefix:        conf.Common.Config.Metrics.Graphite.Prefix,
			Percentiles:   []float64{0.5, 0.75, 0.95, 0.99, 0.999},
		}
		go graphite.WithConfig(graphanaConf)
	}

	var err error
	if conf.Common.Config.HTTPS.Enabled {
		err = fasthttp.ListenAndServeTLS(conf.Common.Address, conf.Common.Config.HTTPS.CertFile, conf.Common.Config.HTTPS.PrivateKey, f(r.Handler))
	} else {
		err = fasthttp.ListenAndServe(conf.Common.Address, f(r.Handler))
	}

	if err != nil {
		conf.Common.Logger.Fatal(err)
	}
}
