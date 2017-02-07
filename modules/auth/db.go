package auth

import (
	"errors"

	"github.com/go-xorm/xorm"
)

func sync(e *xorm.Engine) {
	e.Sync2(new(token))
}

type TokenRepo struct {
	Orm *xorm.Engine
}

func (r *TokenRepo) All() ([]token, error) {
	var ts []token
	err := r.Orm.Find(&ts)
	return ts, err
}
func (r *TokenRepo) Remove(t string) error {
	_, err := r.Orm.Where("token =?", t).Delete(new(token))
	return err
}

var errNotFound = errors.New("Not FOUND")

func (r *TokenRepo) Get(tId string) (*token, error) {
	var t token
	has, err := r.Orm.Where("token =?", tId).Get(&t)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errNotFound
	}
	return &t, nil
}
func (r *TokenRepo) Create(t token) error {
	_, err := r.Orm.Insert(t)
	return err
}
