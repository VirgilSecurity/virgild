package coreapi

import metrics "github.com/rcrowley/go-metrics"

type cacheManager struct {
	logger Logger
	cache  RawCache
}

func (m cacheManager) Get(key string, val interface{}) bool {
	var has bool
	var err error

	t := metrics.GetOrRegisterTimer("cache.get", nil)
	t.Time(func() {
		has, err = m.cache.Get(key, val)
	})

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
		return false
	}
	return has
}

func (m cacheManager) Set(key string, val interface{}) {
	var err error

	t := metrics.GetOrRegisterTimer("cache.set", nil)
	t.Time(func() {
		err = m.cache.Set(key, val)
	})

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
	}
}
func (m cacheManager) Del(key string) {
	var err error

	t := metrics.GetOrRegisterTimer("cache.del", nil)
	t.Time(func() {
		err = m.cache.Del(key)
	})

	if err != nil {
		m.logger.Err("Cache Manager: %+v", err)
	}
}
