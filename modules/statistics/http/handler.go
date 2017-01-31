package http

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/valyala/fasthttp"
)

type StatisticRepoSearch interface {
	Search(from, to int64, group core.StatisticDayGroup, token string) ([]core.StatisticGroup, error)
}

type TimeGroup int

const (
	Day         TimeGroup = iota
	Last3Days   TimeGroup = iota
	LastMonth   TimeGroup = iota
	Last3Months TimeGroup = iota
	Year        TimeGroup = iota
	All         TimeGroup = iota
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func GetStatistic(searcher StatisticRepoSearch, logger Logger) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		args := ctx.QueryArgs()
		var from time.Time
		to := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour) // end of the current day

		var group core.StatisticDayGroup
		g := string(args.Peek("group"))
		switch g {
		case "month":
			group = core.Hour
			from = to.AddDate(0, -1, 0)
		case "last_3_months":
			group = core.Day
			from = to.AddDate(0, -3, 0)
		case "year":
			group = core.Month
			from = to.AddDate(-1, 0, 0)
		case "all":
			group = core.Month
			from = to.AddDate(-100, 0, 0)
		case "last_3_days":
			from = to.AddDate(0, 0, -3)
			group = core.Hour
		default:
			from = to.AddDate(0, 0, -1)
			group = core.Hour
		}

		token := string(args.Peek("token"))
		resp, err := searcher.Search(from.Unix(), to.Unix(), group, token)
		if err != nil {
			logger.Printf("Get statistic info(group=%v token=%v): %+v", group, token, err)
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(resp)
		if err != nil {
			logger.Printf("Cannot marshal result: %+v", err)
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Success("application/json", b)
	}
}

type LastActionsRepo interface {
	Get(until int64) ([]core.RequestStatistics, error)
}

func LastActions(repo LastActionsRepo, logger Logger) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var until int64
		suntil := string(ctx.QueryArgs().Peek("until"))
		if suntil != "" {
			var err error
			until, err = strconv.ParseInt(suntil, 10, 64)
			if err != nil {
				ctx.Error("until parameter must be int", fasthttp.StatusBadRequest)
				return
			}
		}

		r, err := repo.Get(until)
		if err != nil {
			logger.Printf("Last actions statistic info: %+v", err)
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(r)
		if err != nil {
			logger.Printf("Cannot marshal result: %+v", err)
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Success("application/json", b)
	}
}
