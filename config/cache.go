package config

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"

	"github.com/allegro/bigcache"
	"github.com/dchest/siphash"
	"github.com/pkg/errors"
)

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

type Logger interface {
	Printf(format string, args ...interface{})
}

type CacheManager struct {
	Cache  *bigcache.BigCache
	Logger Logger
}

func (m *CacheManager) Get(key string, val interface{}) bool {
	r, err := m.Cache.Get(key)
	if _, ok := err.(*bigcache.EntryNotFoundError); ok {
		return false
	}

	if err != nil {
		m.Logger.Printf("Cache manager get(key=%v) internal error: %+v", key, errors.WithStack(err))
		return false
	}

	if len(r) == 0 {
		return false
	}
	err = json.Unmarshal(r, val)
	if err != nil {
		m.Logger.Printf("Cache manager: get(%v) unmarshal error: %+v", key, errors.WithStack(err))
		return false
	}

	return true
}

func (m *CacheManager) Set(key string, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		m.Logger.Printf("Cache manager: set(%v) marshal error: %+v", key, errors.WithStack(err))
		return
	}
	err = m.Cache.Set(key, b)
	if err != nil {
		m.Logger.Printf("Cache manager: set(%v,%s) internal error: %+v", key, b, errors.WithStack(err))
		return
	}
}

func (m *CacheManager) Del(key string) {
	_, err := m.Cache.Get(key)
	if err != nil {
		m.Logger.Printf("Cache manager: del(%v) internal error: %+v", key, errors.WithStack(err))
		return
	}
	err = m.Cache.Set(key, nil)
	if err != nil {
		m.Logger.Printf("Cache manager: del(%v,nil) internal error: %+v", key, errors.WithStack(err))
	}
}
