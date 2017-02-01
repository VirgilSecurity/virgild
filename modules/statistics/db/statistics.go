package db

import (
	"fmt"
	"time"

	"github.com/VirgilSecurity/virgild/modules/statistics/core"
	"github.com/go-xorm/xorm"
)

func Sync(orm *xorm.Engine) error {
	return orm.Sync2(new(core.RequestStatistics))
}

type StatisticRepository struct {
	Orm *xorm.Engine
}

func (r *StatisticRepository) Add(s core.RequestStatistics) error {
	_, err := r.Orm.InsertOne(s)
	return err
}

func (r *StatisticRepository) Search(from, to int64, group core.StatisticDayGroup, token string) ([]core.StatisticGroup, error) {
	var result []core.StatisticGroup
	q := r.Orm.Where("date >=?", from).And("date <?", to)
	if token != "" {
		q = q.And("token = ?", token)
	}

	dGroup := setGroupByDate(group)
	err := q.GroupBy(fmt.Sprintf(`token, %v,
	case
		when method = 'GET' and resource like '/v4/card/%%' then 1
		when method = 'DELETE' and resource like '/v4/card/%%' then 2
		when method = 'POST' and resource = '/v4/card/actions/search' then 3
		when method = 'POST' and resource = '/v4/card' then 4
		else 0
	end `, dGroup)).
		Select(fmt.Sprintf(`%v as date, count(*) as count,token,
	case
		when method = 'GET' and resource like '/v4/card/%%' then 1
		when method = 'DELETE' and resource like '/v4/card/%%' then 2
		when method = 'POST' and resource = '/v4/card/actions/search' then 3
		when method = 'POST' and resource = '/v4/card' then 4
		else 0
	end as endpoint`, dGroup)).
		Table(new(core.RequestStatistics)).
		Find(&result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func setGroupByDate(group core.StatisticDayGroup) string {
	switch group {
	case core.Hour:
		round := int64(time.Hour / time.Second)
		return fmt.Sprintf("date/%v*%v", round, round)
	case core.Day:
		round := int64(24 * time.Hour / time.Second)
		return fmt.Sprintf("date/%v*%v", round, round)
	default: // Month
		return "date_month"
	}
}

func (r *StatisticRepository) Get(until int64) ([]core.RequestStatistics, error) {
	var result []core.RequestStatistics
	q := r.Orm.Limit(50)
	if until != 0 {
		q.Where("id<?", until)
	}
	err := q.OrderBy("id desc").Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
