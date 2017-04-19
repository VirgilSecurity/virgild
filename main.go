package main

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card"
	"github.com/VirgilSecurity/virgild/modules/healthcheck"
	_ "github.com/VirgilSecurity/virgild/plugins/cache"
	_ "github.com/VirgilSecurity/virgild/plugins/logs"
	_ "github.com/VirgilSecurity/virgild/plugins/metrics"
	"github.com/namsral/flag"
	metrics "github.com/rcrowley/go-metrics"
)

var (
	address         string
	httpsEnabled    bool
	httpsSertificat string
	httpsPrivateKey string
)

func init() {
	flag.StringVar(&address, "address", ":8080", "Address of service")
	flag.StringVar(&httpsSertificat, "https-certificate", "", "The path of the certificate file")
	flag.StringVar(&httpsPrivateKey, "https-private-key", "", "The path of private key file")
	flag.BoolVar(&httpsEnabled, "https-enabled", false, "Enable HTTPS mode")
}

func main() {
	flag.Parse()

	c := coreapi.Init()
	card.Init(c)
	healthcheck.Init(c)

	c.Common.Logger.Info("Start listening...")

	h := responseMetrics{c.HTTP.Router}
	var err error
	if httpsEnabled {
		err = http.ListenAndServeTLS(address, httpsSertificat, httpsPrivateKey, h)
	} else {
		err = http.ListenAndServe(address, h)
	}

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
