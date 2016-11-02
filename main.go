package main

import (
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/database"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/http"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/local"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/remote"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/sync"
	"log"
	"os"
)

func main() {
	orm := database.MakeDatabase("sqlite3:test.db")
	sr := remote.MakeRemoteStorage("AT.690efbee018f626722658e1a660df013f9c0c18b21edbf845ab1d52cfbee499f", remote.RemoteConfig{
		CardsServiceAddress:         "https://cards-stg.virgilsecurity.com",
		ReadonlyCardsServiceAddress: "https://cards-ro-stg.virgilsecurity.com",
		AppID: "d32b745ec2f3ab47add5d89a18f41f5076dc93ccfb5f3c6a575aef58506a24ec",
		PublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEA8jJqWY5hm4tvmnM6QXFdFCErRCnoYdhVNjFggffSCoc=
-----END PUBLIC KEY-----`),
	})

	l := log.New(os.Stderr, "", log.LstdFlags)
	sl := local.Local{
		Repo: &database.CardRepository{
			Orm: orm,
		},
	}

	router := http.MakeRouter(&controllers.Controller{
		Storage: sync.Sync{
			Local:  sl,
			Remote: sr,
			Logger: l,
		},
	}, l)

	fasthttp.ListenAndServe(":8081", router.GetHandleRequest())
}
