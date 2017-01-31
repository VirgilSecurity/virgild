package statistics

import (
	"strings"
	"time"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/VirgilSecurity/virgild/modules/statistics/db"
	"github.com/VirgilSecurity/virgild/modules/statistics/http"
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
		Middleware:   makeStatisticsMiddleware(repo, conf.Common.Logger),
		GetStatistic: http.GetStatistic(repo, conf.Common.Logger),
		LastActions:  http.LastActions(repo, conf.Common.Logger),
	}
}

type statisticsRepository interface {
	Add(s core.RequestStatistics) error
}

type logger interface {
	Printf(format string, args ...interface{})
}

func makeStatisticsMiddleware(repo statisticsRepository, logger logger) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			next(ctx)

			if ctx.Response.StatusCode() != fasthttp.StatusOK {
				return
			}

			token := string(ctx.Request.Header.Peek("Authorization"))
			i := strings.Index(token, ` `)
			if i != -1 {
				token = token[i:]
			}
			d := time.Now().UTC()
			dm := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()) // first day of month
			err := repo.Add(core.RequestStatistics{
				Date:      d.Unix(),
				DateMonth: dm.Unix(),
				Token:     token,
				Method:    string(ctx.Method()),
				Resource:  string(ctx.Path()),
			})

			if err != nil {
				logger.Printf("Statistic internal error: %+v", err)
			}
		}

	}
}
