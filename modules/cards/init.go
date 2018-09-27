package cards

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/VirgilSecurity/virgild/modules/cards/db"
	"github.com/VirgilSecurity/virgild/modules/cards/http"
	"github.com/VirgilSecurity/virgild/modules/cards/middleware"
	"github.com/VirgilSecurity/virgild/modules/cards/mode"
	"github.com/VirgilSecurity/virgild/modules/cards/validator"
	"github.com/valyala/fasthttp"
	virgil "gopkg.in/virgil.v4"
)

type cardMode interface {
	Get(id string) (*core.Card, error)
	Search(c *virgil.Criteria) ([]core.Card, error)
	Create(req *core.CreateCardRequest) (*core.Card, error)
	Revoke(req *core.RevokeCardRequest) error
}

type CardsHandlers struct {
	GetCard     fasthttp.RequestHandler
	SearchCards fasthttp.RequestHandler
	CreateCard  fasthttp.RequestHandler
	RevokeCard  fasthttp.RequestHandler
	CountCards  fasthttp.RequestHandler
}

func Init(conf *config.App) *CardsHandlers {
	err := db.Sync(conf.Common.DB)
	if err != nil {
		conf.Common.Logger.Fatalln("Cannot sync db", err)
	}

	respWrap := http.MakeResponseWrapper(conf.Common.Logger)
	getCard, searchCards, createCard, revokeCard := makeCardMode(conf)

	if conf.Cards.Mode != config.CardModeCache {
		if conf.Cards.Signer.Card != nil { // first run
			card := conf.Cards.Signer.Card
			createCard(&core.CreateCardRequest{
				Info: virgil.CardModel{
					Identity:     card.Identity,
					IdentityType: card.IdentityType,
					Scope:        card.Scope,
					DeviceInfo:   card.DeviceInfo,
					Data:         card.Data,
				},
				Request: virgil.SignableRequest{
					Snapshot: card.Snapshot,
					Meta: virgil.RequestMeta{
						Signatures: card.Signatures,
					},
				},
			})
		}

		signer := middleware.MakeSigner(conf.Cards.Signer.CardID, conf.Cards.Signer.PrivateKey)
		createCard = middleware.SignCreateRequest(signer, createCard)
		revokeCard = middleware.SignRevokeRequest(signer, revokeCard)
	}

	return &CardsHandlers{
		GetCard:     respWrap(http.GetCard(getCard)),
		SearchCards: respWrap(http.SearchCards(middleware.SetApplicationScopForSearch(validator.SearchCards(searchCards)))),
		CreateCard:  respWrap(http.CreateCard(validator.CreateCard(createCard))),
		RevokeCard:  respWrap(http.RevokeCard(validator.RevokeCard(revokeCard))),
		CountCards:  respWrap(http.GetCountCards(&db.CardRepository{Orm: conf.Common.DB})),
	}

}

func makeCardMode(conf *config.App) (get core.GetCard, search core.SearchCards, create core.CreateCard, revoke core.RevokeCard) {
	cardRepo := &db.MetricsCardRepository{R: db.CardRepository{
		Orm:   conf.Common.DB,
		Cache: conf.Cards.Remote.Cache,
	}}

	switch conf.Cards.Mode {
	case config.CardModeCache:
		remote := mode.RemoteCardsMiddleware{
			Client: conf.Cards.Remote.VClient,
		}
		get = remote.Get
		search = remote.Search
		create = remote.Create
		revoke = remote.Revoke
	case config.CardModeLocal:
		dummy := mode.DummyCardsMiddleware{}
		local := mode.LocalCardsMiddleware{
			Repo: cardRepo,
		}
		get = local.Get(dummy.Get)
		search = local.Search(dummy.Search)
		create = local.Create(dummy.Create)
		revoke = local.Revoke(dummy.Revoke)
	case config.CardModeSync:
		remote := mode.RemoteCardsMiddleware{
			Client: conf.Cards.Remote.VClient,
		}
		local := mode.LocalCardsMiddleware{
			Repo: cardRepo,
		}

		get = local.Get(remote.Get)
		search = local.Search(remote.Search)
		create = local.Create(remote.Create)
		revoke = local.Revoke(remote.Revoke)
	default:
		conf.Common.Logger.Fatalf("Unsupported cards mode (%v)\n", conf.Cards.Mode)
		return nil, nil, nil, nil
	}

	cache := mode.CacheCardsMiddleware{Manager: conf.Common.Cache}
	return cache.Get(get),
		cache.Search(search),
		cache.Create(create),
		cache.Revoke(revoke)
}
