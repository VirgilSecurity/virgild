package middleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type fakeStatisticRepo struct {
	mock.Mock
}

func (f *fakeStatisticRepo) Add(s core.RequestStatistics) error {
	args := f.Called(s)
	return args.Error(0)
}

type fakeLogger struct {
	mock.Mock
}

func (f *fakeLogger) Printf(format string, args ...interface{}) {
	f.Called()
}

func TestMakeStatisticsMiddleware_StatusCodeNotOk_Skeep(t *testing.T) {
	ctx := fasthttp.RequestCtx{}
	next := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	repo := new(fakeStatisticRepo)

	s := MakeStatisticsMiddleware(repo, nil)
	s(next)(&ctx)

	repo.AssertNotCalled(t, "Add", mock.Anything)
}

func TestMakeStatisticsMiddleware_RepoReturnErr_LogErr(t *testing.T) {
	ctx := fasthttp.RequestCtx{}
	next := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
	repo := new(fakeStatisticRepo)
	repo.On("Add", mock.Anything).Return(fmt.Errorf("error"))
	log := new(fakeLogger)
	log.On("Printf").Once()

	s := MakeStatisticsMiddleware(repo, log)
	s(next)(&ctx)

	log.AssertExpectations(t)
}

func TestMakeStatisticsMiddleware_Add(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("foobar.com")
	req.Header.SetMethod("GET")
	req.SetRequestURI("/foo/bar/baz")
	req.Header.Add("Authorization", "bearer 123")
	ctx.Init(&req, nil, nil)

	next := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
	repo := new(fakeStatisticRepo)
	repo.On("Add", mock.MatchedBy(func(stat core.RequestStatistics) bool {
		now := time.Now().UTC()
		return stat.DateMonth == time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Unix() &&
			stat.Method == "GET" &&
			stat.Resource == "/foo/bar/baz" &&
			stat.Token == "123"
	})).Return(nil).Once()

	s := MakeStatisticsMiddleware(repo, nil)
	s(next)(&ctx)

	repo.AssertExpectations(t)
}
