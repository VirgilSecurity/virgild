package cards

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/VirgilSecurity/virgild/modules/cards/db"
	"github.com/VirgilSecurity/virgild/modules/cards/http"
	"github.com/VirgilSecurity/virgild/modules/cards/middleware"
	"github.com/VirgilSecurity/virgild/modules/cards/mode"
	"github.com/VirgilSecurity/virgild/modules/cards/utils"
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
	mode := makeCardMode(conf)

	createCard, revokeCard := mode.Create, mode.Revoke

	if conf.Cards.Mode != config.CardModeCache {
		signer := middleware.MakeSigner(conf.Cards.Signer.CardID, conf.Cards.Signer.PrivateKey)
		if conf.Cards.Signer.Card != nil { // first run
			card := conf.Cards.Signer.Card
			mode.Create(&core.CreateCardRequest{
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

		createCard = middleware.SignCreateRequest(signer, createCard)
		revokeCard = middleware.SignRevokeRequest(signer, revokeCard)
	}

	return &CardsHandlers{
		GetCard:     respWrap(http.GetCard(mode.Get)),
		SearchCards: respWrap(http.SearchCards(middleware.SetApplicationScopForSearch(validator.SearchCards(mode.Search)))),
		CreateCard:  respWrap(http.CreateCard(validator.CreateCard(createCard))),
		RevokeCard:  respWrap(http.RevokeCard(validator.RevokeCard(revokeCard))),
		CountCards:  respWrap(http.GetCountCards(&db.CardRepository{Orm: conf.Common.DB})),
	}

}

func makeCardMode(conf *config.App) cardMode {
	cardRepo := &db.CardRepository{
		Orm:   conf.Common.DB,
		Cache: conf.Cards.Remote.Cache,
	}
	switch conf.Cards.Mode {
	case config.CardModeCache:
		return &mode.CacheModeHandler{
			Remote: conf.Cards.Remote.VClient,
			Repo:   cardRepo,
		}
	case config.CardModeLocal:
		return &mode.DefaultModeCardHandler{
			Repo: cardRepo,
			Fingerprint: &utils.Fingerprint{
				Crypto: virgil.Crypto(),
			},
		}
	case config.CardModeSync:
		return &mode.AppModeCardHandler{
			Repo:   cardRepo,
			Remote: conf.Cards.Remote.VClient,
		}
	default:
		conf.Common.Logger.Fatalln("Unsupported cards mode (%v)", conf.Cards.Mode)
		return nil
	}
}
