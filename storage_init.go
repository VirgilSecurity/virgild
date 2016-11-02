package main

import (
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/local"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/remote"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/sync"
)

func MakeStorage() controllers.Storage {
	switch config.Mode {
	case "local":
		return makeLocalStorage()
	case "remote":
		return makeRemoteStorage()
	case "sync":
		return sync.Sync{
			Local:  makeLocalStorage(),
			Remote: makeRemoteStorage(),
			Logger: MakeLogger(),
		}
	}
	MakeLogger().Fatal("Unknown mode type")
	return nil
}

func makeLocalStorage() controllers.Storage {
	return local.Local{
		Repo: MakeCardRepository(),
	}
}

func makeRemoteStorage() controllers.Storage {
	conf := remote.RemoteConfig{
		CardsServiceAddress:         config.RemoteService.CardsServiceAddress,
		ReadonlyCardsServiceAddress: config.RemoteService.ReadonlyCardsServiceAddress,
		PublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEA5Fle51URZN2seVuToVQKSFZ8OkF051jlUjBuM9OZSHk=
-----END PUBLIC KEY-----`),
		AppID: "d32b745ec2f3ab47add5d89a18f41f5076dc93ccfb5f3c6a575aef58506a24ec",
	}
	return remote.MakeRemoteStorage(config.RemoteService.Token, conf)
}
