package plugin_cache

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/coocood/freecache"
	"github.com/dchest/siphash"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
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

	return freeCache{
		Cache:         freecache.NewCache(cacheSize * 1024 * 1024),
		ExpireSeconds: int(cacheDuration / time.Second),
		Hasher:        h,
	}, nil
}

type hasher interface {
	Sum64(string) int64
}

type freeCache struct {
	Cache         *freecache.Cache
	Hasher        hasher
	ExpireSeconds int
}

func (m freeCache) Get(key string, val interface{}) (bool, error) {
	hash := m.Hasher.Sum64(key)
	r, err := m.Cache.GetInt(hash)
	if err == freecache.ErrNotFound {
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

func (m freeCache) Set(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "Cache: set(%v) marshal error", key)
	}
	hash := m.Hasher.Sum64(key)
	err = m.Cache.SetInt(hash, b, m.ExpireSeconds)
	if err != nil {
		return errors.Wrapf(err, "Cache: set(%v,%s) internal error", key, b)
	}
	return nil
}

func (m freeCache) Del(key string) error {
	hash := m.Hasher.Sum64(key)
	m.Cache.DelInt(hash)

	return nil
}

type siph struct {
	k0, k1 uint64
}

func (h siph) Sum64(data string) int64 {
	return int64(siphash.Hash(h.k0, h.k1, []byte(data)))
}
func newHasher() (hasher, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	return siph{binary.BigEndian.Uint64(key), binary.BigEndian.Uint64(key[8:])}, nil
}
