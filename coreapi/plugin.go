package coreapi

var (
	loggers       map[string]func() (Logger, error)
	cachers       map[string]func() (RawCache, error)
	metricsOutput map[string]func() error
)

func init() {
	loggers = make(map[string]func() (Logger, error))
	cachers = make(map[string]func() (RawCache, error))
	metricsOutput = make(map[string]func() error)
}

func RegisterLogger(key string, makeF func() (Logger, error)) {
	loggers[key] = makeF
}

func RegisterCache(key string, makeF func() (RawCache, error)) {
	cachers[key] = makeF
}

func RegisterMetrics(key string, makeF func() error) {
	metricsOutput[key] = makeF
}

type RawCache interface {
	Get(key string, val interface{}) (bool, error)
	Set(key string, val interface{}) error
	Del(key string) error
}
