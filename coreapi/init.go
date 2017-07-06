package coreapi

import (
	"log"
	"os"

	"github.com/bmizerany/pat"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	loggerType string
	cacheType  string
)

func init() {
	flag.StringVar(&loggerType, "logger-type", "file", "Logger type")
	flag.StringVar(&cacheType, "cache-type", "mem", "Cache type")
}

func Init() Core {
	loggerF, ok := loggers[loggerType]
	if !ok {
		log.Fatalln("Core.init: Logger type (", loggerType, ") are not registred")
	}
	l, err := loggerF()
	if err != nil {
		log.Fatalln("Core.init: Cannot create logger:", err)
	}

	cacheF, ok := cachers[cacheType]
	if !ok {
		l.Err("Core.init: Cache type (%s) are not registred", cacheType)
		os.Exit(-1)
	}
	cache, err := cacheF()
	if err != nil {
		l.Err("Core.init: Cannot create cache: %+v", err)
		os.Exit(-1)
	}

	router := pat.New()
	router.Get("/service/metrics", promhttp.Handler())

	app := Core{
		Common: Common{
			Logger: l,
			Cache: &cacheManager{
				logger: l,
				cache:  cache,
			},
		},
		HTTP: HTTP{
			Router:         router,
			WrapAPIHandler: wrapAPIHandler(l),
		},
	}

	return app
}
