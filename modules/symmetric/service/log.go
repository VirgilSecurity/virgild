package service

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/valyala/fasthttp"
)

type SymmetricHandler func(core.SymmetricRepo) core.Response

func Log(s LogStorage, r core.SymmetricRepo) func(f SymmetricHandler) core.Response {
	return func(f SymmetricHandler) core.Response {
		return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
			owner := "unknow"
			r := &LogSymmetricOperation{r, s, owner}

			h := f(r)
			return h(ctx)
		}
	}
}

type LogStorage interface {
	Add(core.LogSymmetricKey) error
}

type LogSymmetricOperation struct {
	r     core.SymmetricRepo
	s     LogStorage
	owner string
}

func (l *LogSymmetricOperation) Create(k core.SymmetricKey) error {
	err := l.r.Create(k)
	if err == nil {
		err = l.s.Add(core.LogSymmetricKey{
			UserID:    k.UserID,
			KeyID:     k.KeyID,
			Created:   time.Now().UTC().Unix(),
			Operation: core.OperationCreateKey,
			WhoId:     l.owner,
		})

	}
	return err
}
func (l *LogSymmetricOperation) Remove(keyID, userID string) error {
	err := l.r.Remove(keyID, userID)
	if err == nil {
		err = l.s.Add(core.LogSymmetricKey{
			UserID:    userID,
			KeyID:     keyID,
			Created:   time.Now().UTC().Unix(),
			Operation: core.OperationRemoveKey,
			WhoId:     l.owner,
		})

	}
	return err
}
func (l *LogSymmetricOperation) Get(keyID, userID string) (k *core.SymmetricKey, err error) {
	k, err = l.r.Get(keyID, userID)
	if err == nil {
		err = l.s.Add(core.LogSymmetricKey{
			UserID:    userID,
			KeyID:     keyID,
			Created:   time.Now().UTC().Unix(),
			Operation: core.OperationGetKey,
			WhoId:     k.UserID,
		})
	}
	return
}
