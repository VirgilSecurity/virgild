package main

import (
	"github.com/go-xorm/xorm"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/database"
)

var orm *xorm.Engine

func getOrm() *xorm.Engine {
	if orm == nil {

		orm = database.MakeDatabase(config.DatabseConnection)
	}
	return orm
}
func MakeCardRepository() *database.CardRepository {
	return &database.CardRepository{
		Orm: getOrm(),
	}
}
