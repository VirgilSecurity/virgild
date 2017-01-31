package db

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/go-xorm/xorm"
)

func Sync(orm *xorm.Engine) error {
	return orm.Sync2(new(core.SqlCard))
}

type CardRepository struct {
	Cache time.Duration
	Orm   *xorm.Engine
}

func (r *CardRepository) Get(id string) (*core.SqlCard, error) {
	cs := new(core.SqlCard)
	has, err := r.Orm.Where("card_id =?", id).Get(cs)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, core.ErrorEntityNotFound
	}
	return cs, nil
}

func (r *CardRepository) Find(identitis []string, identityType string, scope string) ([]core.SqlCard, error) {
	var cs []core.SqlCard
	q := r.Orm.In("identity", identitis).
		And("scope = ?", scope)

	if identityType != "" {
		q = q.And("identity_type = ?", identityType)
	}
	err := q.Find(&cs)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (r *CardRepository) Add(cs core.SqlCard) error {
	if cs.Scope == "global" || len(cs.Card) == 0 {
		cs.ExpireAt = time.Now().UTC().Add(r.Cache)
	} else {
		cs.ExpireAt = time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	_, err := r.Orm.InsertOne(cs)
	return err
}

func (r *CardRepository) DeleteById(id string) error {
	_, err := r.Orm.Where("card_id =?", id).Delete(new(core.SqlCard))
	return err
}

func (r *CardRepository) DeleteBySearch(identitis []string, identityType string, scope string) error {
	q := r.Orm.In("identity", identitis).
		And("scope = ?", scope)

	if identityType != "" {
		q = q.And("identity_type = ?", identityType)
	}

	_, err := q.Delete(new(core.SqlCard))
	return err
}

func (r *CardRepository) Count() (int64, error) {
	count, err := r.Orm.Where("error_code!=0").And("expire_at>?", time.Now().UTC()).Count(new([]core.SqlCard))
	if err != nil {
		return 0, err
	}
	return count, nil
}
