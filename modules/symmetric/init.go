package symmetric

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/valyala/fasthttp"
)

type SymmetricModule struct {
	GetKey         fasthttp.RequestHandler
	CreateKey      fasthttp.RequestHandler
	GetKeysForUser fasthttp.RequestHandler
	GetUsersForKey fasthttp.RequestHandler
}

func Init(conf *config.App) *SymmetricModule {
	Sync(conf.Common.DB)
	repo := &SymmetricKeyRepo{
		Orm: conf.Common.DB,
	}
	wrap := MakeResponseWrapper(conf.Common.Logger)
	return &SymmetricModule{
		GetKey:         wrap(getKey(repo)),
		CreateKey:      wrap(createKey(repo)),
		GetKeysForUser: wrap(getKeysByUser(repo)),
		GetUsersForKey: wrap(getUsersByKey(repo)),
	}
}
