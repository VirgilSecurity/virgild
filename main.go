package main

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card"
	_ "github.com/VirgilSecurity/virgild/plugins/cache"
	_ "github.com/VirgilSecurity/virgild/plugins/logs"
	_ "github.com/VirgilSecurity/virgild/plugins/metrics"
	"github.com/namsral/flag"
	metrics "github.com/rcrowley/go-metrics"
)

var address string

func init() {
	flag.StringVar(&address, "address", ":8080", "Address of service")
}

func main() {
	flag.Parse()

	c := coreapi.Init()
	card.Init(c)

	c.Common.Logger.Info("Start listening...")
	err := http.ListenAndServe(address, responseMetrics{c.HTTP.Router})
	if err != nil {
		c.Common.Logger.Err("HTTP server return err: %v", err)
	}
}

type responseMetrics struct {
	router http.Handler
}

func (m responseMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := metrics.GetOrRegisterTimer("response", nil)
	t.Time(func() {
		m.router.ServeHTTP(w, r)
	})
}
