package app

import (
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/local"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/remote"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/sync"
	"os"
)

func makeStorage() controllers.Storage {
	switch config.Mode {
	case "local":
		return makeLocalStorage()
	case "sync":
		return &sync.Sync{
			Local:  makeLocalStorage(),
			Remote: makeRemoteStorage(),
		}
	}
	logger.Fatal("Unknown mode type.")
	os.Exit(2)
	return nil
}

func makeLocalStorage() controllers.Storage {
	return &local.Local{
		Repo:   makeCardRepository(),
		Logger: logger,
	}
}

func makeRemoteStorage() controllers.Storage {
	conf := remote.RemoteConfig{
		CardsServiceAddress:         config.RemoteService.CardsServiceAddress,
		ReadonlyCardsServiceAddress: config.RemoteService.ReadonlyCardsServiceAddress,
	}
	return remote.MakeRemoteStorage(config.RemoteService.Token, conf)
}
