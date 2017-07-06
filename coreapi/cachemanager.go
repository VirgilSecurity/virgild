package coreapi

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheManagerMetric = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "duration_seconds",
		Subsystem:  "cache_manager",
		Namespace:  "virgild",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(cacheManagerMetric)
}

type cacheManager struct {
	logger Logger
	cache  RawCache
}

func (m cacheManager) Get(key string, val interface{}) bool {
	var has bool
	var err error

	t := prometheus.NewTimer(cacheManagerMetric.WithLabelValues("get"))
	has, err = m.cache.Get(key, val)
	t.ObserveDuration()

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
		return false
	}
	return has
}

func (m cacheManager) Set(key string, val interface{}) {
	var err error

	t := prometheus.NewTimer(cacheManagerMetric.WithLabelValues("set"))
	err = m.cache.Set(key, val)
	t.ObserveDuration()

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
	}
}
func (m cacheManager) Del(key string) {
	var err error

	t := prometheus.NewTimer(cacheManagerMetric.WithLabelValues("del"))
	err = m.cache.Del(key)
	t.ObserveDuration()

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
	}
}
