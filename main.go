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
	Help:      "HTTP handler latency in seconds",
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

	http.Handle("/", corsHandler(httpDuration(c.HTTP.Router)))
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
func corsHandler(hander http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "X-Virgil-Request-Id, X-Virgil-Request-Sign, X-Virgil-Response-Id, X-Virgil-Response-Sign, X-Virgil-Access-Token, X-Virgil-Application-Token, X-Virgil-Request-Uuid, X-Virgil-Request-Sign-Virgil-Card-ID, X-Virgil-Request-Sign-Pk-Id, X-Virgil-Authentication, Content-Type, User-Agent, Origin, Authorization, Accept, DNT, X-Requested-With, If-Modified-Since, Cache-Control")
		w.Header().Set("Access-Control-Expose-Headers", "X-Virgil-Response-Id, X-Virgil-Response-Sign")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		hander.ServeHTTP(w, r)
	})
}
