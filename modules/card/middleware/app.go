package middleware

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/pkg/errors"
)

type AppStore interface {
	GetById(id string) (*core.Application, error)
}

type TokenStore interface {
	GetByValue(val string) (*core.Token, error)
}

type AppMiddleware struct {
	Cache      coreapi.Cache
	AppStore   AppStore
	TokenStore TokenStore
}

func (m *AppMiddleware) RequestApp(next coreapi.APIHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		var appID string
		owner := core.GetOwnerRequest(req.Context())
		if owner == "" {
			return next(req)
		}
		has := m.Cache.Get(owner, &appID)
		if !has {
			token, err := m.TokenStore.GetByValue(owner)
			if err != nil {
				return nil, errors.Wrapf(err, "AppMiddleware.GetTokenByValue(%s)", owner)
			}
			app, _ := m.AppStore.GetById(token.Application)
			appID = app.ID
			m.Cache.Set(owner, appID)
		}
		ctx := core.SetOwnerRequest(req.Context(), appID)
		return next(req.WithContext(ctx))
	}
}
