package main

import (
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func Sync(orm *xorm.Engine) error {
	return orm.Sync2(new(cardSql))
}

type ImpSqlCardRepository struct {
	Cache time.Duration
	Orm   *xorm.Engine
}

func (r *ImpSqlCardRepository) Get(id string) (*cardSql, error) {
	cs := new(cardSql)
	has, err := r.Orm.Where("card_id =?", id).Get(cs)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, ErrorEntityNotFound
	}
	return cs, nil
}

func (r *ImpSqlCardRepository) Find(identitis []string, identityType string, scope string) ([]cardSql, error) {
	var cs []cardSql
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

func (r *ImpSqlCardRepository) Add(cs cardSql) error {
	if cs.Scope == "global" || len(cs.Card) == 0 {
		cs.ExpireAt = time.Now().UTC().Add(r.Cache)
	} else {
		cs.ExpireAt = time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	_, err := r.Orm.InsertOne(cs)
	return err
}

func (r *ImpSqlCardRepository) DeleteById(id string) error {
	_, err := r.Orm.Where("card_id =?", id).Delete(new(cardSql))
	return err
}

func (r *ImpSqlCardRepository) DeleteBySearch(identitis []string, identityType string, scope string) error {
	q := r.Orm.In("identity", identitis).
		And("scope = ?", scope)

	if identityType != "" {
		q = q.And("identity_type = ?", identityType)
	}

	_, err := q.Delete(new(cardSql))
	return err
}
