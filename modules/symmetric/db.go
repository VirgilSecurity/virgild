package symmetric

import "github.com/go-xorm/xorm"

func Sync(e *xorm.Engine) error {
	return e.Sync2(new(SymmetricKey))
}

type SymmetricKey struct {
	KeyID        string `xorm:"key_id" json:"key_id"`
	UserID       string `xorm:"user_id" json:"user_id"`
	EncryptedKey []byte `xorm:"encrypted_key" json:"encrypted_key"`
}

type SymmetricKeyRepo struct {
	Orm *xorm.Engine
}

func (r *SymmetricKeyRepo) Create(k SymmetricKey) error {
	_, err := r.Orm.InsertOne(k)
	return err
}

func (r *SymmetricKeyRepo) Remove(keyID, userID string) error {
	_, err := r.Orm.Where("key_id = ?", keyID).And("user_id = ?", userID).Delete(new(SymmetricKey))
	return err
}

func (r *SymmetricKeyRepo) Get(keyID, userID string) (*SymmetricKey, error) {
	var k SymmetricKey
	has, err := r.Orm.Where("key_id = ?", keyID).And("user_id = ?", userID).Get(&k)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrorEntityNotFound
	}
	return &k, nil
}

func (r *SymmetricKeyRepo) KeysByUser(userID string) (ks []SymmetricKey, err error) {
	err = r.Orm.Where("user_id = ?", userID).Find(&ks)
	return
}

func (r *SymmetricKeyRepo) UsersByKey(keyID string) (ks []SymmetricKey, err error) {
	err = r.Orm.Where("key_id = ?", keyID).Find(&ks)
	return
}
