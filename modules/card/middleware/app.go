package middleware

import (
	"fmt"
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

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
		has := m.Cache.Get(owner, &appID)
		if !has {
			token, _ := m.TokenStore.GetByValue(owner)
			// if err != nil {
			// 	return next(req)
			// }
			app, _ := m.AppStore.GetById(token.Application)
			appID = app.ID
			m.Cache.Set(owner, appID)

			fmt.Println("Find app ", appID, "for", owner)
		}
		ctx := core.SetOwnerRequest(req.Context(), appID)
		return next(req.WithContext(ctx))
	}
}
