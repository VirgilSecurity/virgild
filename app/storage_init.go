package app

import (
	"github.com/virgilsecurity/virgild/controllers"
	"github.com/virgilsecurity/virgild/storage/local"
	"github.com/virgilsecurity/virgild/storage/remote"
	"github.com/virgilsecurity/virgild/storage/sync"
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
	return remote.MakeRemoteStorage(config.RemoteService.Token, logger, conf)
}
