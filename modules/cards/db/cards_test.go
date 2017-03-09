// +build integration

package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
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

func TestGet_CardNotExist_ReturnErr(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	repo := CardRepository{Orm: orm}
	_, err = repo.Get("1")

	assert.Equal(t, core.ErrorEntityNotFound, err)
}

func TestGet_CardExist_ReturnVal(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	orm.Insert(core.SqlCard{CardID: "1"})
	repo := CardRepository{Orm: orm}
	c, _ := repo.Get("1")

	assert.NotNil(t, c)
}

func TestFind_CardsNotExist_ReturnEmpty(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	repo := CardRepository{Orm: orm}
	cs, _ := repo.Find([]string{"test"}, "", "")
	assert.Len(t, cs, 0)
}

func TestFind_CardsExist_ReturnVal(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	orm.Insert(core.SqlCard{
		Identity:     "test1",
		Scope:        "global",
		IdentityType: "nick",
		ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}, core.SqlCard{
		Identity:     "test2",
		Scope:        "app",
		IdentityType: "nick",
		ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}, core.SqlCard{
		Identity:     "test3",
		Scope:        "app",
		IdentityType: "email",
		ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	})
	repo := CardRepository{Orm: orm}

	type Data struct {
		Identity     []string
		Scope        string
		IdentityType string
		Result       []core.SqlCard
	}
	table := []Data{
		// filter by identities
		Data{[]string{"test2", "test3"}, "app", "", []core.SqlCard{core.SqlCard{
			Identity:     "test2",
			Scope:        "app",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}, core.SqlCard{
			Identity:     "test3",
			Scope:        "app",
			IdentityType: "email",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
		// filter by scope
		Data{[]string{"test1", "test2"}, "global", "", []core.SqlCard{core.SqlCard{
			Identity:     "test1",
			Scope:        "global",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
		// filter by identity type
		Data{[]string{"test1", "test2", "test3"}, "app", "email", []core.SqlCard{core.SqlCard{
			Identity:     "test3",
			Scope:        "app",
			IdentityType: "email",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
	}

	for _, v := range table {
		cs, _ := repo.Find(v.Identity, v.IdentityType, v.Scope)
		assert.EqualValues(t, v.Result, cs, fmt.Sprintln("Identity:", v.Identity, "Type:", v.IdentityType, "Scope:", v.Scope))
	}
}

func TestAdd_ScopeGlobal_SetCorrectCache(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	c := core.SqlCard{
		Identity:     "testScopeGlobal",
		Scope:        "global",
		IdentityType: "nick",
		Card:         []byte("test"),
	}
	repo := CardRepository{time.Hour, orm}
	before := time.Now().Add(time.Hour)
	time.Sleep(time.Second)

	err = repo.Add(c)
	assert.Nil(t, err)

	time.Sleep(time.Second)
	after := time.Now().Add(time.Hour)

	var actual core.SqlCard
	orm.Where("identity=?", "testScopeGlobal").Get(&actual)

	actualExp := time.Unix(actual.ExpireAt, 0)
	assert.True(t, after.After(actualExp))
	assert.True(t, before.Before(actualExp))
}

func TestAdd_CardEmpty_SetCorrectCache(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	c := core.SqlCard{
		Identity:     "testCardEmpty",
		Scope:        "app",
		IdentityType: "nick",
	}
	repo := CardRepository{time.Hour, orm}
	before := time.Now().Add(time.Hour)
	time.Sleep(time.Second)

	err = repo.Add(c)
	assert.Nil(t, err)

	time.Sleep(time.Second)
	after := time.Now().Add(time.Hour)

	var actual core.SqlCard
	orm.Where("identity=?", "testCardEmpty").Get(&actual)

	actualExp := time.Unix(actual.ExpireAt, 0)
	assert.True(t, after.After(actualExp))
	assert.True(t, before.Before(actualExp))
}

func TestAdd_ScopeApplication_NeverExp(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	c := core.SqlCard{
		Identity:     "testScopeApp",
		Scope:        "app",
		IdentityType: "nick",
		Card:         []byte("test"),
	}
	repo := CardRepository{time.Hour, orm}
	err = repo.Add(c)
	assert.Nil(t, err)

	var actual core.SqlCard
	orm.Where("identity=?", "testScopeApp").Get(&actual)

	assert.Equal(t, time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), actual.ExpireAt)
}

func TestMarkDeletedById_Remove(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	orm.InsertOne(core.SqlCard{
		CardID:       "removeByID",
		Identity:     "removeByIDIdenityt",
		Scope:        "removeByIDScope",
		IdentityType: "removeByIDNick",
		ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	})
	repo := CardRepository{time.Hour, orm}
	err = repo.MarkDeletedById("removeByID")
	assert.Nil(t, err)
	var actual core.SqlCard
	has, _ := orm.Where("card_id=?", "removeByID").Get(&actual)

	assert.True(t, has)
	assert.True(t, actual.Deleted)
}

func TestDeleteById_Remove(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	orm.InsertOne(core.SqlCard{
		CardID:       "removeByID",
		Identity:     "removeByIDIdenityt",
		Scope:        "removeByIDScope",
		IdentityType: "removeByIDNick",
		ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	})
	repo := CardRepository{time.Hour, orm}
	err = repo.DeleteById("removeByID")
	assert.Nil(t, err)
	var actual core.SqlCard
	has, _ := orm.Where("card_id=?", "removeByID").Get(&actual)

	assert.False(t, has)
}

func TestDeleteBySearch_Remove(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	type Data struct {
		Identity     []string
		Scope        string
		IdentityType string
		Result       []core.SqlCard
	}

	table := []Data{
		// filter by identities
		Data{[]string{"test2", "test3"}, "app", "", []core.SqlCard{core.SqlCard{
			Identity:     "test1",
			Scope:        "global",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
		// filter by scope
		Data{[]string{"test1", "test2"}, "global", "", []core.SqlCard{core.SqlCard{
			Identity:     "test2",
			Scope:        "app",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}, core.SqlCard{
			Identity:     "test3",
			Scope:        "app",
			IdentityType: "email",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
		// filter by identity type
		Data{[]string{"test1", "test2", "test3"}, "app", "email", []core.SqlCard{core.SqlCard{
			Identity:     "test1",
			Scope:        "global",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}, core.SqlCard{
			Identity:     "test2",
			Scope:        "app",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}}},
	}
	for _, v := range table {
		orm.Where("1=1").Delete(new(core.SqlCard))

		orm.Insert(core.SqlCard{
			Identity:     "test1",
			Scope:        "global",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}, core.SqlCard{
			Identity:     "test2",
			Scope:        "app",
			IdentityType: "nick",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		}, core.SqlCard{
			Identity:     "test3",
			Scope:        "app",
			IdentityType: "email",
			ExpireAt:     time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		})

		repo := CardRepository{Orm: orm}
		err := repo.DeleteBySearch(v.Identity, v.IdentityType, v.Scope)
		assert.Nil(t, err)

		var cs []core.SqlCard
		orm.Find(&cs)
		assert.EqualValues(t, v.Result, cs, fmt.Sprintln("Identity:", v.Identity, "Type:", v.IdentityType, "Scope:", v.Scope))
	}
}

func TestCount_ReturnCount(t *testing.T) {
	orm, err := initDB()
	assert.Nil(t, err, "Cannot init db")
	defer finDB(orm)

	orm.Insert(core.SqlCard{ // Skeeped because ExpireAt < now
		Identity:     "test2",
		Scope:        "app",
		IdentityType: "nick",
		ExpireAt:     time.Now().AddDate(-10, 0, 0).Unix(),
	}, core.SqlCard{
		Identity:     "test3",
		Scope:        "app",
		IdentityType: "email",
		ExpireAt:     time.Now().AddDate(10, 0, 0).Unix(),
	}, core.SqlCard{
		Identity:     "test4",
		Scope:        "app",
		IdentityType: "email",
		ExpireAt:     time.Now().AddDate(0, 1, 0).Unix(),
	})

	repo := CardRepository{Orm: orm}
	count, err := repo.Count()

	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
}
