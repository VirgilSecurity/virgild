//// +build integration

package db

import (
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func initDB() (*xorm.Engine, error) {
	orm, err := xorm.NewEngine("sqlite3", "")
	if err != nil {
		return nil, fmt.Errorf("Cannot create db: %v", err)
	}
	err = Sync(orm)
	if err != nil {
		return nil, fmt.Errorf("Cannot migrate sql schema: %v", err)
	}
	return orm, nil
}

func finDB(orm *xorm.Engine) {
	orm.Close()
}

func TestGet_KeyNotExist_ReturnErr(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	repo := SymmetricKeyRepo{Orm: orm}
	_, err = repo.Get("1", "1")

	assert.Equal(t, core.ErrorEntityNotFound, err)
}

func TestGet_ReturnVal(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	expected := &core.SymmetricKey{UserID: "test", KeyID: "1234", EncryptedKey: []byte("encrypted key")}
	orm.InsertOne(expected)

	repo := SymmetricKeyRepo{Orm: orm}
	actual, _ := repo.Get("1234", "test")

	assert.Equal(t, expected, actual)
}

func TestKeysByUser_ReturnVal(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	expected := core.SymmetricKey{UserID: "test", KeyID: "1234", EncryptedKey: []byte("encrypted key")}
	orm.Insert(expected,
		&core.SymmetricKey{UserID: "test1", KeyID: "1234", EncryptedKey: []byte("encrypted key")},
		&core.SymmetricKey{UserID: "test2", KeyID: "1234", EncryptedKey: []byte("encrypted key")})

	repo := SymmetricKeyRepo{Orm: orm}
	actual, _ := repo.KeysByUser("test")

	assert.Len(t, actual, 1)
	assert.Equal(t, core.KeyUserPair{UserID: expected.UserID, KeyID: expected.KeyID}, actual[0])
}

func TestUsersByKey_ReturnVal(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	expected := core.SymmetricKey{UserID: "test", KeyID: "1234", EncryptedKey: []byte("encrypted key")}
	orm.Insert(expected,
		&core.SymmetricKey{UserID: "test", KeyID: "666", EncryptedKey: []byte("encrypted key")},
		&core.SymmetricKey{UserID: "test", KeyID: "555", EncryptedKey: []byte("encrypted key")})

	repo := SymmetricKeyRepo{Orm: orm}
	actual, _ := repo.UsersByKey("1234")

	assert.Len(t, actual, 1)
	assert.Equal(t, core.KeyUserPair{UserID: expected.UserID, KeyID: expected.KeyID}, actual[0])
}
