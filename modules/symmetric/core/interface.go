package core

import "github.com/valyala/fasthttp"

type Response func(ctx *fasthttp.RequestCtx) (interface{}, error)

type SymmetricRepo interface {
	Create(k SymmetricKey) error
	Remove(keyID, userID string) error
	Get(keyID, userID string) (k *SymmetricKey, err error)
}

type ListSymmetricRepo interface {
	KeysByUser(userID string) (ks []KeyUserPair, err error)
	UsersByKey(keyID string) (ks []KeyUserPair, err error)
}
