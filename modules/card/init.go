package card

import (
	"github.com/VirgilSecurity/virgild/coreapi"
	vhttp "github.com/VirgilSecurity/virgild/modules/card/http"
	"github.com/VirgilSecurity/virgild/modules/card/middleware"
	"github.com/VirgilSecurity/virgild/modules/card/validator"
	"github.com/namsral/flag"
)

var (
	raService    string
	cardsService string
)

func init() {
	flag.StringVar(&raService, "card-raservice", "https://ra.virgilsecurity.com", "Addres of Registration authority")
	flag.StringVar(&cardsService, "card-cardsservice", "https://cards.virgilsecurity.com", "Addres of Cards")
}

func Init(c coreapi.Core) {
	apiWrap := c.HTTP.WrapAPIHandler
	cache := cacheCardMiddleware{cache: c.Common.Cache}
	cloud := cloudCard{
		CardsService: cardsService,
		RAService:    raService,
	}

	hGet := middleware.RequestOwner(vhttp.GetCard(cache.GetCard(cloud.getCard)))
	hSearch := middleware.RequestOwner(vhttp.SearchCards(middleware.SetApplicationScopForSearch(validator.SearchCards(cache.SearchCards(cloud.searchCards)))))
	hCreateCard := middleware.RequestOwner(vhttp.CreateCard(validator.CreateCard(cache.CreateCard(cloud.createCard))))
	hRevokeCard := middleware.RequestOwner(vhttp.RevokeCard(validator.RevokeCard(cache.RevokeCard(cloud.revokeCard))))
	hCreateRelation := middleware.RequestOwner(vhttp.CreateRelation(cache.CreateRelations(cloud.createRelation)))
	hRevokeRelation := middleware.RequestOwner(vhttp.RevokeRelation(cache.RevokeRelations(cloud.revokeRelation)))

	r := c.HTTP.Router
	r.Post("/v1/card", apiWrap(hCreateCard))
	r.Del("/v1/card/:id", apiWrap(hRevokeCard))
	r.Post("/v4/card", apiWrap(hCreateCard))
	r.Post("/v4/card/actions/search", apiWrap(hSearch))
	r.Del("/v4/card/:id", apiWrap(hRevokeCard))
	r.Get("/v4/card/:id", apiWrap(hGet))
	r.Post("/v4/card/:id/collections/relations", apiWrap(hCreateRelation))
	r.Del("/v4/card/:id/collections/relations", apiWrap(hRevokeRelation))
}
