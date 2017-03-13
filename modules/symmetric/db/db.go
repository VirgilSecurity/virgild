package db

import (
	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/go-xorm/xorm"
)

func Sync(e *xorm.Engine) error {
	return e.Sync2(new(core.SymmetricKey), new(core.LogSymmetricKey))
}

type LogSymmetricKeyRepo struct {
	Orm *xorm.Engine
}

func (r *LogSymmetricKeyRepo) Add(o core.LogSymmetricKey) error {
	_, err := r.Orm.Insert(o)
	return err
}

type SymmetricKeyRepo struct {
	Orm *xorm.Engine
}

func (r *SymmetricKeyRepo) Create(k core.SymmetricKey) error {
	_, err := r.Orm.InsertOne(k)
	return err
}

func (r *SymmetricKeyRepo) Remove(keyID, userID string) error {
	_, err := r.Orm.Where("key_id = ?", keyID).And("user_id = ?", userID).Delete(new(core.SymmetricKey))
	return err
}

func (r *SymmetricKeyRepo) Get(keyID, userID string) (*core.SymmetricKey, error) {
	var k core.SymmetricKey
	has, err := r.Orm.Where("key_id = ?", keyID).And("user_id = ?", userID).Get(&k)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, core.ErrorEntityNotFound
	}
	return &k, nil
}

func (r *SymmetricKeyRepo) KeysByUser(userID string) (ks []core.KeyUserPair, err error) {
	err = r.Orm.Where("user_id = ?", userID).Table(new(core.SymmetricKey)).Find(&ks)
	return
}

func (r *SymmetricKeyRepo) UsersByKey(keyID string) (ks []core.KeyUserPair, err error) {
	err = r.Orm.Where("key_id = ?", keyID).Table(new(core.SymmetricKey)).Find(&ks)
	return
}
