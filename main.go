package main

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card"
	"github.com/VirgilSecurity/virgild/modules/healthcheck"
	_ "github.com/VirgilSecurity/virgild/plugins/cache"
	_ "github.com/VirgilSecurity/virgild/plugins/logs"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	address         string
	httpsEnabled    bool
	httpsSertificat string
	httpsPrivateKey string
)
var rpc = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:      "duration_seconds",
	Subsystem: "http",
	Namespace: "virgild",
})

func init() {
	flag.StringVar(&address, "address", ":8080", "Address of service")
	flag.BoolVar(&httpsEnabled, "https-enabled", false, "Enable HTTPS mode")
	flag.StringVar(&httpsSertificat, "https-certificate", "", "The path of the certificate file")
	flag.StringVar(&httpsPrivateKey, "https-private-key", "", "The path of private key file")

	prometheus.MustRegister(rpc)
}

func main() {
	flag.Parse()

	c := coreapi.Init()
	card.Init(c)
	healthcheck.Init(c)

	c.Common.Logger.Info("Start listening address %v ...", address)

	http.Handle("/", httpDuration(c.HTTP.Router))
	var err error
	if httpsEnabled {
		err = http.ListenAndServeTLS(address, httpsSertificat, httpsPrivateKey, nil)
	} else {
		err = http.ListenAndServe(address, nil)
	}

	if err != nil {
		c.Common.Logger.Err("HTTP server return err: %v", err)
	}
}

func httpDuration(hander http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := prometheus.NewTimer(rpc)
		hander.ServeHTTP(w, r)
		t.ObserveDuration()
	})
}
