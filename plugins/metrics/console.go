package plugin_metrics

import (
	"log"
	"os"
	"time"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/namsral/flag"
	metrics "github.com/rcrowley/go-metrics"
)

var consoleInterval time.Duration

func init() {
	flag.DurationVar(&consoleInterval, "metrics-console-interval", time.Minute, "Interval between flushing data to console output")
	coreapi.RegisterMetrics("console", initConsoleMetrics)
}

func initConsoleMetrics() error {
	go metrics.LogScaled(metrics.DefaultRegistry, consoleInterval, time.Microsecond, log.New(os.Stdout, "[METRICS] ", log.LstdFlags|log.LUTC))
	return nil
}
