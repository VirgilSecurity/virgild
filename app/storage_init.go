package app

import (
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/local"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/remote"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/sync"
)

func makeStorage() controllers.Storage {
	switch config.Mode {
	case "local":
		return makeLocalStorage()
	case "remote":
		return makeRemoteStorage()
	case "sync":
		return &sync.Sync{
			Local:  makeLocalStorage(),
			Remote: makeRemoteStorage(),
		}
	}
	makeLogger().Fatal("Unknown mode type")
	return nil
}

func makeLocalStorage() controllers.Storage {
	return &local.Local{
		Repo:   makeCardRepository(),
		Logger: makeLogger(),
	}
}

func makeRemoteStorage() controllers.Storage {
	conf := remote.RemoteConfig{
		CardsServiceAddress:         config.RemoteService.CardsServiceAddress,
		ReadonlyCardsServiceAddress: config.RemoteService.ReadonlyCardsServiceAddress,
	}
	return remote.MakeRemoteStorage(config.RemoteService.Token, conf)
}
