package plugin_cache

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/allegro/bigcache"
	"github.com/dchest/siphash"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
)

var (
	cacheDuration time.Duration
	cacheSize     int
)

func init() {
	flag.DurationVar(&cacheDuration, "cache-duration", time.Hour, "Cache duration")
	flag.IntVar(&cacheSize, "cache-size", 1024, "Cache size (mb)")

	coreapi.RegisterCache("mem", makeMemoryCache)

}

func makeMemoryCache() (coreapi.RawCache, error) {
	h, err := newHasher()
	if err != nil {
		return nil, errors.Wrap(err, "Create hash function")
	}
	bconf := bigcache.DefaultConfig(cacheDuration)
	bconf.HardMaxCacheSize = cacheSize
	bconf.Hasher = h
	bc, err := bigcache.NewBigCache(bconf)
	if err != nil {
		return nil, errors.Wrap(err, "Create memory cache")
	}

	metrics.NewFunctionalGauge(func() int64 {
		return int64(bc.Len())
	})

	return cache{bc}, nil
}

type cache struct {
	Cache *bigcache.BigCache
}

func (m cache) Get(key string, val interface{}) (bool, error) {
	r, err := m.Cache.Get(key)
	if _, ok := err.(*bigcache.EntryNotFoundError); ok {
		return false, nil
	}

	if err != nil {
		return false, errors.Wrapf(err, "Cache: get(key=%v) internal error", key)
	}

	if len(r) == 0 {
		return false, nil
	}
	err = json.Unmarshal(r, val)
	if err != nil {
		return false, errors.Wrapf(err, "Cache: get(%v) unmarshal error", key)
	}

	return true, nil
}

func (m cache) Set(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "Cache: set(%v) marshal error", key)
	}
	err = m.Cache.Set(key, b)
	if err != nil {
		return errors.Wrapf(err, "Cache: set(%v,%s) internal error", key, b)
	}
	return nil
}

func (m cache) Del(key string) error {
	_, err := m.Cache.Get(key)
	if _, ok := err.(*bigcache.EntryNotFoundError); ok {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "Cache: del(%v) get key", key)
	}
	err = m.Cache.Set(key, nil)
	if err != nil {
		return errors.Wrapf(err, "Cache: del(%v) set nil", key)
	}
	return nil
}

type siph struct {
	k0, k1 uint64
}

func (h siph) Sum64(data string) uint64 {
	return siphash.Hash(h.k0, h.k1, []byte(data))
}
func newHasher() (bigcache.Hasher, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	return siph{binary.BigEndian.Uint64(key), binary.BigEndian.Uint64(key[8:])}, nil
}
