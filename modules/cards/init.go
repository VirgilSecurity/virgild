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
	db.Sync(conf.Common.DB)

	respWrap := http.MakeResponseWrapper(conf.Common.Logger)
	mode := makeCardMode(conf)
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

	return &CardsHandlers{
		GetCard:     respWrap(http.GetCard(mode.Get)),
		SearchCards: respWrap(http.SearchCards(middleware.SetApplicationScopForSearch(validator.SearchCards(mode.Search)))),
		CreateCard:  respWrap(http.CreateCard(validator.CreateCard(middleware.SignCreateRequest(signer, mode.Create)))),
		RevokeCard:  respWrap(http.RevokeCard(validator.RevokeCard(middleware.SignRevokeRequest(signer, mode.Revoke)))),
		CountCards:  respWrap(http.GetCountCards(&db.CardRepository{Orm: conf.Common.DB})),
	}

}

func makeCardMode(conf *config.App) cardMode {
	cardRepo := &db.CardRepository{
		Orm: conf.Common.DB,
	}

	if conf.Cards.Remote != nil {
		return &mode.AppModeCardHandler{
			Repo:   cardRepo,
			Remote: conf.Cards.Remote.VClient,
		}
	}
	return &mode.DefaultModeCardHandler{
		Repo: cardRepo,
		Fingerprint: &utils.Fingerprint{
			Crypto: virgil.Crypto(),
		},
	}

}
