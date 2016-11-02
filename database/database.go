package database

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"strings"
)

func MakeDatabase(db string) *xorm.Engine {
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

	err = engine.Sync(new(sqlmodels.CardSql))
	if err != nil {
		panic("Cannot sync db")
	}

	return engine
}
