package card

import (
	"time"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/background"
	"github.com/VirgilSecurity/virgild/modules/card/db"
	vhttp "github.com/VirgilSecurity/virgild/modules/card/http"
	"github.com/VirgilSecurity/virgild/modules/card/integrations"
	"github.com/VirgilSecurity/virgild/modules/card/middleware"
	"github.com/VirgilSecurity/virgild/modules/card/validator"
	"github.com/namsral/flag"
	"github.com/rubenv/sql-migrate"
)

var (
	raService         string
	cardsService      string
	devPortalService  string
	devPortalLogin    string
	devPortalPassword string
)

func init() {
	flag.StringVar(&raService, "card-raservice", "https://ra.virgilsecurity.com", "Addres of Registration authority")
	flag.StringVar(&cardsService, "card-cardsservice", "https://cards.virgilsecurity.com", "Addres of Cards")

	flag.StringVar(&devPortalService, "card-dev-portal-service", "https://devportal.virgilsecurity.com", "Addres of Development portal")
	flag.StringVar(&devPortalLogin, "card-dev-portal-login", "", "Login for Development portal")
	flag.StringVar(&devPortalPassword, "card-dev-portal-password", "", "Password for Development portal")
}

func Init(c coreapi.Core) {

	applyMigrations(c)

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

	if devPortalLogin != "" && devPortalPassword != "" {
		devPortalClient := &inegrations.DevPortalClient{
			Address: devPortalService,
		}

		err := devPortalClient.Authorize(devPortalLogin, devPortalPassword)
		if err != nil {
			c.Common.Logger.Err("Cannot authorize to dev portal (login: %s password: %s): %+v", devPortalLogin, devPortalPassword, err)
			c.Common.Logger.Warn("synchronization is disabled")
		} else {
			tokenStore := &db.TokenStore{DB: c.Common.DB}
			appStore := &db.ApplicationsStore{DB: c.Common.DB}

			c.Scheduler.Add(time.Minute, background.UpdateAppsCronJob(appStore, devPortalClient), "UpdateApplications")
			c.Scheduler.Add(time.Minute, background.UpdateTokensCronJob(tokenStore, devPortalClient), "UpdateTokens")

			appMidware := middleware.AppMiddleware{
				Cache:      c.Common.Cache,
				TokenStore: tokenStore,
				AppStore:   appStore,
			}
			hGet = appMidware.RequestApp(hGet)
			hSearch = appMidware.RequestApp(hSearch)
			hCreateCard = appMidware.RequestApp(hCreateCard)
			hRevokeCard = appMidware.RequestApp(hRevokeCard)
			hCreateRelation = appMidware.RequestApp(hCreateRelation)
			hRevokeRelation = appMidware.RequestApp(hRevokeRelation)
		}

	}

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

func applyMigrations(c coreapi.Core) {
	migrations := &migrate.AssetMigrationSource{
		Asset:    db.Asset,
		AssetDir: db.AssetDir,
		Dir:      "db/migrations",
	}

	_, err := migrate.Exec(c.Common.DB.DB, c.Common.DB.DriverName(), migrations, migrate.Up)
	if err != nil {
		c.Common.Logger.Err("Card module: apply migrations: %+v", err)
		panic(nil)
	}
}
