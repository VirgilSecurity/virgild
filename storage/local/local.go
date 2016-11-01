package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func MakeLocalStorage(db string) Local {
	i := strings.Index(db, ":")
	if i == -1 {
		panic("Database connection should be a format {driver}:{dataSourceName}")
	}
	driver := db[:i]
	s := db[i+1:]
	engine, err := xorm.NewEngine(driver, s)
	if err != nil {
		panic("Cannot connect to db")
	}

	err = engine.Sync(new(CardSql))
	if err != nil {
		panic("Cannot sync with db")
	}

	return Local{
		engine: engine,
	}
}

type CardSql struct {
	Id           string `xorm:"PK"`
	Identity     string `xorm:"Index"`
	IdentityType string
	Scope        string
	Card         string
}

type Local struct {
	engine *xorm.Engine
}

func (s Local) GetCard(id string) (models.CardResponse, error) {
	var r models.CardResponse
	var c CardSql
	has, err := s.engine.Where("id = ?", id).Get(&c)
	if err != nil {
		return r, err
	}
	if !has {
		return r, errors.New("Card was not found")
	}
	err = json.Unmarshal([]byte(c.Card), &r)
	return r, err
}

func (s Local) SearchCards(c models.Criteria) (models.CardsResponse, error) {
	var r models.CardsResponse
	var cs []CardSql
	q := s.engine.In("identity", c.Identities).And("scope = ?", c.Scope)
	if c.IdentityType != "" {
		q = q.And("identity_type = ?", c.IdentityType)
	}
	err := q.Find(&cs)
	if err != nil {
		return r, err
	}
	for _, v := range cs {
		var cr models.CardResponse
		err = json.Unmarshal([]byte(v.Card), &cr)
		if err != nil {
			return r, err
		}
		r = append(r, cr)
	}
	return r, err
}

func (s Local) CreateCard(c models.CardResponse) (models.CardResponse, error) {
	var cr models.CardRequest
	err := json.Unmarshal(c.Snapshot, &cr)
	if err != nil {
		return c, err
	}

	jCard, err := json.Marshal(c)
	if err != nil {
		return c, err
	}
	cs := CardSql{
		Id:           c.ID,
		Identity:     cr.Identity,
		IdentityType: cr.IdentityType,
		Scope:        cr.Scope,
		Card:         string(jCard[:]),
	}
	_, err = s.engine.Insert(cs)
	if err != nil {
		fmt.Println("Insert errore:", err)
	}
	return c, nil
}

func (s Local) RevokeCard(id string, c models.CardResponse) error {
	_, err := s.engine.Id(id).Delete(new(CardSql))
	return err
}
