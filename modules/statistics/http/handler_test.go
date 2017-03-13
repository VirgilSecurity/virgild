package http

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type fakeLogger struct {
	mock.Mock
}

func (f *fakeLogger) Printf(format string, args ...interface{}) {
	f.Called()
}

type fakeStatisticRepo struct {
	mock.Mock
}

func (f *fakeStatisticRepo) Search(from, to int64, group core.StatisticDayGroup, token string) (r []core.StatisticGroup, err error) {
	args := f.Called(from, to, group, token)
	r, _ = args.Get(0).([]core.StatisticGroup)
	err = args.Error(1)
	return
}

func (f *fakeStatisticRepo) Get(until int64, count int) (r []core.RequestStatistics, err error) {
	args := f.Called(until, count)
	r, _ = args.Get(0).([]core.RequestStatistics)
	err = args.Error(1)
	return
}

func TestGetStatistic_RepoReturnErr_ReturnErr(t *testing.T) {
	var ctx fasthttp.RequestCtx

	repo := new(fakeStatisticRepo)
	repo.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Error"))

	l := new(fakeLogger)
	l.On("Printf").Once()

	s := GetStatistic(repo, l)
	s(&ctx)

	l.AssertExpectations(t)
	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
}

func TestGetStatistic_ParsParam_ReturnResult(t *testing.T) {
	type param struct {
		Token string
		Quert string
		From  int64
		Group core.StatisticDayGroup
	}
	to := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour) // end of the current day

	table := []param{
		param{
			Token: "",
			Quert: "",
			From:  to.AddDate(0, 0, -1).Unix(),
			Group: core.Hour,
		},
		param{
			Token: "123",
			Quert: "?token=123&group=month",
			From:  to.AddDate(0, -1, 0).Unix(),
			Group: core.Day,
		},
		param{
			Token: "123",
			Quert: "?token=123&group=last_3_months",
			From:  to.AddDate(0, -3, 0).Unix(),
			Group: core.Day,
		},
		param{
			Token: "123",
			Quert: "?token=123&group=year",
			From:  to.AddDate(-1, 0, 0).Unix(),
			Group: core.Month,
		},
		param{
			Token: "123",
			Quert: "?token=123&group=all",
			From:  to.AddDate(-100, 0, 0).Unix(),
			Group: core.Month,
		},
		param{
			Token: "123",
			Quert: "?token=123&group=last_3_days",
			From:  to.AddDate(0, 0, -3).Unix(),
			Group: core.Hour,
		},
	}
	for _, v := range table {
		expected := []core.StatisticGroup{
			core.StatisticGroup{
				Date:     12,
				Token:    "1234",
				Endpoint: 1,
				Count:    1,
			},
		}
		var ctx fasthttp.RequestCtx
		ctx.Request.SetRequestURI("/test" + v.Quert)
		repo := new(fakeStatisticRepo)
		repo.On("Search", v.From, to.Unix(), v.Group, v.Token).Return(expected, nil)

		s := GetStatistic(repo, nil)
		s(&ctx)

		var actual []core.StatisticGroup
		err := json.Unmarshal(ctx.Response.Body(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
		assert.Equal(t, expected, actual)
	}
}

func TestLastActions_QueryInvalid_Return404(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("/test?until=asd")
	s := LastActions(nil, nil)
	s(&ctx)

	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode())
}

func TestLastActions_RepoReturnErr_Return500AndLogIt(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("/test?until=123")

	repo := new(fakeStatisticRepo)
	repo.On("Get", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Error"))

	l := new(fakeLogger)
	l.On("Printf").Once()

	s := LastActions(repo, l)
	s(&ctx)

	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
	l.AssertExpectations(t)
}

func TestLastActions_ReturnResult(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("/test?until=123")

	expected := []core.RequestStatistics{
		core.RequestStatistics{
			Date:     12,
			Method:   "GET",
			Resource: "/1234",
			Token:    "123",
		},
	}
	repo := new(fakeStatisticRepo)
	repo.On("Get", int64(123), COUNT_LAST_ACTIONS).Return(expected, nil)

	s := LastActions(repo, nil)
	s(&ctx)

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))

	var actual []core.RequestStatistics
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
