package config

import (
	metrics "github.com/rcrowley/go-metrics"
	virgil "gopkg.in/virgil.v4"
)

type MetricsVirgilClient struct {
	C *virgil.Client
}

func (c *MetricsVirgilClient) GetCard(id string) (card *virgil.Card, err error) {
	t := metrics.GetOrRegisterTimer("cards-service.get", nil)
	t.Time(func() {
		card, err = c.C.GetCard(id)
	})
	return
}
func (c *MetricsVirgilClient) SearchCards(crit *virgil.Criteria) (cards []*virgil.Card, err error) {
	t := metrics.GetOrRegisterTimer("cards-service.search", nil)
	t.Time(func() {
		cards, err = c.C.SearchCards(crit)
	})
	return
}
func (c *MetricsVirgilClient) CreateCard(req *virgil.SignableRequest) (card *virgil.Card, err error) {
	t := metrics.GetOrRegisterTimer("cards-service.create", nil)
	t.Time(func() {
		card, err = c.C.CreateCard(req)
	})
	return
}
func (c *MetricsVirgilClient) RevokeCard(req *virgil.SignableRequest) (err error) {
	t := metrics.GetOrRegisterTimer("cards-service.revoke", nil)
	t.Time(func() {
		err = c.C.RevokeCard(req)
	})
	return
}

type MetricsCacheManager struct {
	C Cache
}

func (m *MetricsCacheManager) Get(key string, v interface{}) (has bool) {
	t := metrics.GetOrRegisterTimer("cache.get", nil)
	t.Time(func() {
		has = m.C.Get(key, v)
	})
	return
}

func (m *MetricsCacheManager) Set(key string, v interface{}) {
	t := metrics.GetOrRegisterTimer("cache.set", nil)
	t.Time(func() {
		m.C.Set(key, v)
	})
}

func (m *MetricsCacheManager) Del(key string) {
	t := metrics.GetOrRegisterTimer("cache.del", nil)
	t.Time(func() {
		m.C.Del(key)
	})
}
