package statistics

import (
	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/statistics/db"
	"github.com/VirgilSecurity/virgild/modules/statistics/http"
	"github.com/VirgilSecurity/virgild/modules/statistics/middleware"
	"github.com/valyala/fasthttp"
)

type StatisticHandlers struct {
	Middleware   func(fasthttp.RequestHandler) fasthttp.RequestHandler
	GetStatistic fasthttp.RequestHandler
	LastActions  fasthttp.RequestHandler
}

func Init(conf *config.App) *StatisticHandlers {
	db.Sync(conf.Common.DB)

	repo := &db.StatisticRepository{
		Orm: conf.Common.DB,
	}

	return &StatisticHandlers{
		Middleware:   middleware.MakeStatisticsMiddleware(repo, conf.Common.Logger),
		GetStatistic: http.GetStatistic(repo, conf.Common.Logger),
		LastActions:  http.LastActions(repo, conf.Common.Logger),
	}
}
