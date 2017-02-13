package symmetric

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/symmetric/db"
	"github.com/VirgilSecurity/virgild/modules/symmetric/handler"
	"github.com/VirgilSecurity/virgild/modules/symmetric/http"
	"github.com/VirgilSecurity/virgild/modules/symmetric/service"
	"github.com/valyala/fasthttp"
)

type SymmetricModule struct {
	GetKey         fasthttp.RequestHandler
	CreateKey      fasthttp.RequestHandler
	GetKeysForUser fasthttp.RequestHandler
	GetUsersForKey fasthttp.RequestHandler
}

func Init(conf *config.App) *SymmetricModule {
	db.Sync(conf.Common.DB)
	repo := &db.SymmetricKeyRepo{
		Orm: conf.Common.DB,
	}
	s := &db.LogSymmetricKeyRepo{Orm: conf.Common.DB}
	l := service.Log(s, repo)
	wrap := http.MakeResponseWrapper(conf.Common.Logger)

	return &SymmetricModule{
		GetKey:         wrap(l(handler.GetKey)),
		CreateKey:      wrap(l(handler.CreateKey)),
		GetKeysForUser: wrap(handler.GetKeysByUser(repo)),
		GetUsersForKey: wrap(handler.GetUsersByKey(repo)),
	}
}
