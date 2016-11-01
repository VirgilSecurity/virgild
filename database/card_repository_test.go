// +build integration

package database

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"testing"
)

func Before() *xorm.Engine {
	orm, err := xorm.NewEngine("sqlite3", "")
	if err != nil {
		panic("Cannot open a database")
	}
	orm.Sync(new(CardSql))
	return orm
}

var fakeData []CardSql = []CardSql{
	CardSql{
		Identity:     "identity1",
		IdentityType: "nick",
		Scope:        "global",
		Id:           "id1",
		Card:         "some information",
	},
	CardSql{
		Identity:     "identity1_a",
		IdentityType: "nick",
		Scope:        "application",
		Id:           "id1_app",
		Card:         "some information",
	},
	CardSql{
		Identity:     "identity2",
		IdentityType: "email",
		Scope:        "global",
		Id:           "id2",
		Card:         "some information",
	},
	CardSql{
		Identity:     "identity2_a",
		IdentityType: "email",
		Scope:        "appplication",
		Id:           "id2_app",
		Card:         "some information",
	},
	CardSql{
		Identity:     "identity3",
		IdentityType: "global",
		Scope:        "global",
		Id:           "id3",
		Card:         "some information",
	},
}

func AddFakeData(orm *xorm.Engine) {
	for _, v := range fakeData {
		orm.Insert(v)
	}

}

func GetById(identities ...string) []CardSql {
	r := make([]CardSql, 0, len(identities))
	for _, v := range fakeData {
		for _, identity := range identities {
			if v.Identity == identity {
				r = append(r, v)
			}
		}
	}
	return r
}

func Test_Get_EmptyResult_ReturnNil(t *testing.T) {
	orm := Before()
	defer orm.Close()

	repo := CardRepository{
		Orm: orm,
	}

	r, err := repo.Get("asdf")

	assert.Nil(t, r)
	assert.Nil(t, err)
}

func Test_Get_ReturnErr_ReturnErr(t *testing.T) {
	orm := Before()
	orm.Close()

	repo := CardRepository{
		Orm: orm,
	}
	_, err := repo.Get("asdf")
	assert.NotNil(t, err)
}

func Test_Get_HasResult_ReturnVal(t *testing.T) {
	orm := Before()
	defer orm.Close()
	card := CardSql{
		Identity:     "some identity",
		IdentityType: "global",
		Scope:        "global",
		Id:           "id",
		Card:         "some information",
	}
	orm.Insert(card)

	repo := CardRepository{
		Orm: orm,
	}

	r, err := repo.Get(card.Id)

	assert.NotNil(t, r)
	assert.Equal(t, card, *r)
	assert.Nil(t, err)
}

func Test_Find_FilledOnlyIdentities_ReturnVal(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	expected := GetById("identity1", "identity2")

	repo := CardRepository{
		Orm: orm,
	}

	actual, err := repo.Find(models.Criteria{
		Identities: []string{"identity1", "identity2"},
	})

	assert.Nil(t, err)

	assert.NotNil(t, actual)
	assert.EqualValues(t, expected, actual)
}

func Test_Search_ReturnErr_ReturnErr(t *testing.T) {
	orm := Before()
	orm.Close()

	repo := CardRepository{
		Orm: orm,
	}
	_, err := repo.Find(models.Criteria{
		Identities: []string{"identity1", "identity2"},
	})
	assert.NotNil(t, err)
}

func Test_Find_FilledIdentitiesScope_ReturnVal(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	expected := GetById("identity1", "identity2")

	repo := CardRepository{
		Orm: orm,
	}

	actual, err := repo.Find(models.Criteria{
		Identities: []string{"identity1", "identity2", "identity2_a"},
		Scope:      "global",
	})

	assert.Nil(t, err)

	assert.NotNil(t, actual)
	assert.EqualValues(t, expected, actual)
}

func Test_Find_FilledIdentitiesIdentityType_ReturnVal(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	expected := GetById("identity2_a", "identity2")

	repo := CardRepository{
		Orm: orm,
	}

	actual, err := repo.Find(models.Criteria{
		Identities:   []string{"identity1", "identity2", "identity2_a"},
		IdentityType: "email",
	})

	assert.Nil(t, err)

	assert.NotNil(t, actual)
	assert.EqualValues(t, expected, actual)
}

func Test_Find_FilledIdentitiesIdentityTypeScope_ReturnVal(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	expected := GetById("identity1")

	repo := CardRepository{
		Orm: orm,
	}

	actual, err := repo.Find(models.Criteria{
		Identities:   []string{"identity1", "identity2", "identity2_a", "identity1_a"},
		IdentityType: "nick",
		Scope:        "global",
	})

	assert.Nil(t, err)

	assert.NotNil(t, actual)
	assert.EqualValues(t, expected, actual)
}

func Test_Find_EmptyResult_ReturnEmptyArra(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	repo := CardRepository{
		Orm: orm,
	}
	actual, err := repo.Find(models.Criteria{
		Identities: []string{"identity1_empty"},
	})
	assert.Nil(t, err)
	assert.Len(t, actual, 0)
}

func Test_Add(t *testing.T) {
	orm := Before()
	defer orm.Close()

	repo := CardRepository{
		Orm: orm,
	}
	expected := CardSql{
		Identity:     "some identity",
		IdentityType: "global",
		Scope:        "global",
		Id:           "id",
		Card:         "some information",
	}
	repo.Add(expected)

	var actual CardSql
	orm.ID("id").Get(&actual)
	assert.Equal(t, expected, actual)
}

func Test_Add_ReturnErr_ReturnErr(t *testing.T) {
	orm := Before()
	orm.Close()

	repo := CardRepository{
		Orm: orm,
	}
	expected := CardSql{
		Identity:     "some identity",
		IdentityType: "global",
		Scope:        "global",
		Id:           "id",
		Card:         "some information",
	}
	err := repo.Add(expected)

	assert.NotNil(t, err)
}

func Test_Delete(t *testing.T) {
	orm := Before()
	defer orm.Close()
	AddFakeData(orm)

	repo := CardRepository{
		Orm: orm,
	}
	repo.Delete("id1")
	var actual CardSql
	has, _ := orm.ID("id").Get(&actual)
	assert.False(t, has)
}

func Test_Delete_ReturnErr_ReturnErr(t *testing.T) {
	orm := Before()
	orm.Close()

	repo := CardRepository{
		Orm: orm,
	}
	repo.Delete("id1")
	var actual CardSql
	_, err := orm.ID("id").Get(&actual)

	assert.NotNil(t, err)
}
