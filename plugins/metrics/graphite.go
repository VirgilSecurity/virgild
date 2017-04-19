package plugin_metrics

import (
	"net"
	"time"

	"github.com/VirgilSecurity/virgild/coreapi"
	graphite "github.com/cyberdelia/go-metrics-graphite"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	metrics "github.com/rcrowley/go-metrics"
)

var (
	graphiteInterval time.Duration
	graphitePrefix   string
	graphiteAddress  string
)

func init() {
	flag.StringVar(&graphiteAddress, "metrics-graphite-address", "", "Address of graphite service where will be sending metrics (if this parameter is empty then metrics will not send)")
	flag.DurationVar(&graphiteInterval, "metrics-graphite-interval", time.Minute, "Interval between flushing data to graphite")
	flag.StringVar(&graphitePrefix, "metrics-graphite-prefix", "", "Prefix for VirgilD in graphite")

	coreapi.RegisterMetrics("graphite", initGraphiteMetrics)
}

func initGraphiteMetrics() error {
	if graphiteAddress == "" {
		return errors.New("Graphite Address is not passed")
	}
	addr, err := net.ResolveTCPAddr("tcp", graphiteAddress)
	if err != nil {
		return errors.Wrapf(err, "Cannot resolve graphite address (%s)", graphiteAddress)
	}

	graphanaConf := graphite.Config{
		Addr:          addr,
		Registry:      metrics.DefaultRegistry,
		FlushInterval: graphiteInterval,
		DurationUnit:  time.Microsecond,
		Prefix:        graphitePrefix,
		Percentiles:   []float64{0.5, 0.75, 0.95, 0.99, 0.999},
	}
	go graphite.WithConfig(graphanaConf)

	return nil
}
