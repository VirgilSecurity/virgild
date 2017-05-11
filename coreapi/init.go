package coreapi

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/goji/httpauth"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namsral/flag"
)

var (
	database         string
	loggerType       string
	cacheType        string
	adminLogin       string
	adminPassword    string
	apiAuthMode      string
	metricsType      string
	statisticEnabled bool
)

func init() {
	flag.StringVar(&loggerType, "logger", "file", "Logger type")
	flag.StringVar(&database, "database", "sqlite3:virgild.db", "Databse")
	flag.StringVar(&cacheType, "cache", "mem", "Cache type")
	flag.StringVar(&metricsType, "metrics", "disabled", "Metrics output type")

	// flag.StringVar(&adminLogin, "admin-login", "admin", "Admin login")
	// flag.StringVar(&adminPassword, "admin-password", "admin", "Admin password")
	// flag.StringVar(&apiAuthMode, "api-auth-mode", "no", "API Auth mode(no,local)")
	// flag.BoolVar(&statisticEnabled, "statistic", false, "Request statistics")
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

	db, err := initDB(database)
	if err != nil {
		l.Err("Core.init: Cannot open database: %+v", err)
		os.Exit(-1)
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

	metricsF, ok := metricsOutput[metricsType]
	if !ok {
		l.Err("Core.init: Mtrics output type (%s) are not registred", metricsType)
		os.Exit(-1)
	}
	err = metricsF()
	if err != nil {
		l.Err("Core.init: Cannot init metrics (%s): %+v", metricsType, err)
		os.Exit(-1)
	}

	app := Core{
		Common: Common{
			Logger: l,
			DB:     db,
			Cache: &cacheManager{
				logger: l,
				cache:  cache,
			},
		},
		HTTP: HTTP{
			Router:         pat.New(),
			WrapAPIHandler: wrapAPIHandler(l),
			AdminAuth:      httpauth.SimpleBasicAuth(adminLogin, adminPassword),
		},
	}
	return app
}

func initDB(database string) (*sqlx.DB, error) {

	i := strings.Index(database, ":")
	if i < 0 {
		return nil, fmt.Errorf("Cannot parse database connectin ({provider:connection})")
	}
	driver, connection := database[:i], database[i+1:]
	db, err := sqlx.Open(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to (driver: %v name: %v) database: %v", driver, connection, err)
	}
	return db, nil
}
