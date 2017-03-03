package middleware

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/allegro/bigcache"
	"github.com/valyala/fasthttp"

	"github.com/dchest/siphash"
)

type siph struct {
	k0, k1 uint64
}

func (h siph) Sum64(data string) uint64 {
	return siphash.Hash(h.k0, h.k1, []byte(data))
}
func NewHasher() (bigcache.Hasher, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	return siph{binary.BigEndian.Uint64(key), binary.BigEndian.Uint64(key[8:])}, nil
}

func MakeCache(cache *bigcache.BigCache) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			key := fmt.Sprintf("%s_%s_%s", ctx.Method(), ctx.Path(), ctx.PostBody())
			r, err := cache.Get(key)
			if err == nil && len(key) == 0 {
				ctx.Success("application/json", r)
				return
			}
			next(ctx)
			if ctx.Response.StatusCode() == fasthttp.StatusOK {
				cache.Set(key, ctx.Response.Body())
			}
		}
	}
}
