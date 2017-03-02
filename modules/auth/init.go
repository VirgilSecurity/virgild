package auth

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/valyala/fasthttp"
)

const (
	PermissionGetCard     string = "get_card"
	PermissionSearchCards string = "search_cards"
	PermissionCreateCard  string = "create_card"
	PermissionRevokeCard  string = "revoke_card"
)

type Middleware func(permission string, next fasthttp.RequestHandler) fasthttp.RequestHandler

type AuthHandler struct {
	Middleware  Middleware
	GetTokens   fasthttp.RequestHandler
	RemoveToken fasthttp.RequestHandler
	CreateToken fasthttp.RequestHandler
	UpdateToken fasthttp.RequestHandler
}

func Init(app *config.App) *AuthHandler {
	err := sync(app.Common.DB)
	if err != nil {
		app.Common.Logger.Fatalln("Cannot sync db", err)
	}
	repo := &TokenRepo{
		Orm: app.Common.DB,
	}
	th := tokenHandler{
		repo: repo,
	}
	var m Middleware
	switch app.Auth.Mode {
	case config.AuthModeNo:
		m = noAuth
	case config.AuthModeLocal:
		m = wrap(app.Common.Logger, app.Auth.TokenType, localScopes(repo))
	}
	return &AuthHandler{
		Middleware:  m,
		GetTokens:   th.All,
		RemoveToken: th.Remove,
		CreateToken: th.Create,
		UpdateToken: th.Update,
	}
}
