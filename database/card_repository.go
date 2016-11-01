package database

import (
	"github.com/go-xorm/xorm"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
)

type CardRepository struct {
	Orm *xorm.Engine
}

func (r *CardRepository) Get(id string) (*CardSql, error) {
	cs := new(CardSql)
	has, err := r.Orm.ID(id).Get(cs)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}
	return cs, nil
}

func (r *CardRepository) Find(c models.Criteria) ([]CardSql, error) {
	var cs []CardSql
	q := r.Orm.In("identity", c.Identities)
	if c.Scope != "" {
		q = q.And("scope = ?", c.Scope)
	}
	if c.IdentityType != "" {
		q = q.And("identity_type = ?", c.IdentityType)
	}
	err := q.Find(&cs)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (r *CardRepository) Add(cs CardSql) error {
	_, err := r.Orm.Insert(cs)
	return err
}

func (r *CardRepository) Delete(id string) error {
	_, err := r.Orm.Id(id).Delete(new(CardSql))
	return err
}
