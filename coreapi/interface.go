package coreapi

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/jmoiron/sqlx"
)

type Core struct {
	Common Common
	HTTP   HTTP
}

type Common struct {
	Logger Logger
	DB     *sqlx.DB
	Cache  Cache
}

type HTTP struct {
	WrapAPIHandler func(fun APIHandler) http.Handler
	Router         *pat.PatternServeMux
	//	AdminProtection Middleware
}

// API declaration
type APIHandler func(req *http.Request) (interface{}, error)
type Middleware func(next http.Handler) http.Handler

type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Err(format string, args ...interface{})
}

type Cache interface {
	Get(key string, val interface{}) bool
	Set(key string, val interface{})
	Del(key string)
}