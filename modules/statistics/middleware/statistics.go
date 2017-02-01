package middleware

import (
	"strings"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/valyala/fasthttp"
)

type statisticsRepository interface {
	Add(s core.RequestStatistics) error
}

type logger interface {
	Printf(format string, args ...interface{})
}

func MakeStatisticsMiddleware(repo statisticsRepository, logger logger) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			next(ctx)

			if ctx.Response.StatusCode() != fasthttp.StatusOK {
				return
			}

			token := string(ctx.Request.Header.Peek("Authorization"))
			i := strings.Index(token, ` `)
			if i != -1 {
				token = token[i+1:]
			}
			d := time.Now().UTC()
			dm := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC) // first day of month
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
