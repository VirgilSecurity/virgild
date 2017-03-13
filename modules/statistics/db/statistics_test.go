// +build integration

package db

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	repo StatisticRepository
	from int64 = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	to   int64 = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
)

func TestMain(m *testing.M) {
	e, err := xorm.NewEngine("sqlite3", "")
	if err != nil {
		fmt.Println("Cannot create in-memory sqlite3 database")
		os.Exit(1)
	}

	Sync(e)
	fillDB(e)

	repo.Orm = e
	os.Exit(m.Run())
}

func fillDB(e *xorm.Engine) {
	e.Insert(core.RequestStatistics{
		Date:      time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
		DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Method:    "GET",
		Resource:  "/v4/card/1234",
		Token:     "123",
	},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 1, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "GET",
			Resource:  "/v4/card/1234",
			Token:     "123",
		},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 1, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "POST",
			Resource:  "/v4/card",
		},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 1, 12, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "POST",
			Resource:  "/v4/card",
		},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 1, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "POST",
			Resource:  "/v4/card/actions/search",
			Token:     "333",
		},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 1, 12, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "DELETE",
			Resource:  "/v4/card/1234",
			Token:     "333",
		},
		core.RequestStatistics{
			Date:      time.Date(2016, 1, 2, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "GET",
			Resource:  "/v4/card/1234",
			Token:     "123",
		})
}

func TestSearch_GroupHour(t *testing.T) {
	expected := []core.StatisticGroup{
		core.StatisticGroup{
			Count:    2,
			Date:     time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
			Token:    "123",
			Endpoint: core.GetCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.SearchCardsEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Date:     time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.CreateCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Date:     time.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.CreateCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.RevokeCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "123",
			Date:     time.Date(2016, 1, 2, 10, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.GetCardEndpoint,
		},
	}
	actual, err := repo.Search(from, to, core.Hour, "")
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestSearch_GroupDay(t *testing.T) {
	expected := []core.StatisticGroup{
		core.StatisticGroup{
			Count:    2,
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Token:    "123",
			Endpoint: core.GetCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.SearchCardsEndpoint,
		},
		core.StatisticGroup{
			Count:    2,
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.CreateCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.RevokeCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "123",
			Date:     time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.GetCardEndpoint,
		},
	}
	actual, err := repo.Search(from, to, core.Day, "")
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestSearch_GroupMonth(t *testing.T) {
	expected := []core.StatisticGroup{
		core.StatisticGroup{
			Count:    3,
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Token:    "123",
			Endpoint: core.GetCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.SearchCardsEndpoint,
		},
		core.StatisticGroup{
			Count:    2,
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.CreateCardEndpoint,
		},
		core.StatisticGroup{
			Count:    1,
			Token:    "333",
			Date:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Endpoint: core.RevokeCardEndpoint,
		},
	}
	actual, err := repo.Search(from, to, core.Month, "")
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestGet_UntilZero(t *testing.T) {
	expected := []core.RequestStatistics{
		core.RequestStatistics{
			Id:        7,
			Date:      time.Date(2016, 1, 2, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "GET",
			Resource:  "/v4/card/1234",
			Token:     "123",
		},
		core.RequestStatistics{
			Id:        6,
			Date:      time.Date(2016, 1, 1, 12, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "DELETE",
			Resource:  "/v4/card/1234",
			Token:     "333",
		},
	}
	actual, err := repo.Get(0, 2)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestGet_Until3(t *testing.T) {
	expected := []core.RequestStatistics{
		core.RequestStatistics{
			Id:        2,
			Date:      time.Date(2016, 1, 1, 10, 10, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "GET",
			Resource:  "/v4/card/1234",
			Token:     "123",
		},
		core.RequestStatistics{
			Id:        1,
			Date:      time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
			DateMonth: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Method:    "GET",
			Resource:  "/v4/card/1234",
			Token:     "123",
		},
	}
	actual, err := repo.Get(3, 2)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}
