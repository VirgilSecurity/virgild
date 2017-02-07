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
	sync(app.Common.DB)
	repo := &TokenRepo{
		Orm: app.Common.DB,
	}
	th := tokenHandler{
		repo: repo,
	}
	return &AuthHandler{
		Middleware: func(permission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
			return cardsWrap(app.Common.Logger, getToken([]byte("VIRGIL"), auth(repo, permission)), next)
		},
		GetTokens:   th.All,
		RemoveToken: th.Remove,
		CreateToken: th.Create,
		UpdateToken: th.Update,
	}
}
