package db

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/rcrowley/go-metrics"
)

type MetricsCardRepository struct {
	R CardRepository
}

func (r *MetricsCardRepository) Get(id string) (card *core.SqlCard, err error) {
	t := metrics.GetOrRegisterTimer("db.cards.get", nil)
	t.Time(func() {
		card, err = r.R.Get(id)
	})
	return
}

func (r *MetricsCardRepository) Find(identitis []string, identityType string, scope string) (cards []core.SqlCard, err error) {
	t := metrics.GetOrRegisterTimer("db.cards.find", nil)
	t.Time(func() {
		cards, err = r.R.Find(identitis, identityType, scope)
	})
	return
}

func (r *MetricsCardRepository) Add(cs core.SqlCard) (err error) {
	t := metrics.GetOrRegisterTimer("db.cards.add", nil)
	t.Time(func() {
		err = r.R.Add(cs)
	})
	return
}

func (r *MetricsCardRepository) MarkDeletedById(id string) (err error) {
	t := metrics.GetOrRegisterTimer("db.cards.mark-deleted", nil)
	t.Time(func() {
		err = r.R.MarkDeletedById(id)
	})
	return
}

func (r *MetricsCardRepository) DeleteById(id string) (err error) {
	t := metrics.GetOrRegisterTimer("db.cards.delete", nil)
	t.Time(func() {
		err = r.R.DeleteById(id)
	})
	return
}

func (r *MetricsCardRepository) DeleteBySearch(identitis []string, identityType string, scope string) (err error) {
	t := metrics.GetOrRegisterTimer("db.cards.delete-group", nil)
	t.Time(func() {
		err = r.R.DeleteBySearch(identitis, identityType, scope)
	})
	return
}

func (r *MetricsCardRepository) Count() (count int64, err error) {
	t := metrics.GetOrRegisterTimer("db.cards.count", nil)
	t.Time(func() {
		count, err = r.R.Count()
	})
	return
}
