package main

import (
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/http"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/storage/remote"
	"log"
	"os"
)

func main() {
	s := remote.MakeRemoteStorage("AT.690efbee018f626722658e1a660df013f9c0c18b21edbf845ab1d52cfbee499f", remote.RemoteConfig{
		CardsServiceAddress:         "https://cards-stg.virgilsecurity.com",
		ReadonlyCardsServiceAddress: "https://cards-ro-stg.virgilsecurity.com",
	})

	router := http.MakeRouter(&controllers.Controller{
		Storage: s,
	}, log.New(os.Stderr, "", log.LstdFlags))

	fasthttp.ListenAndServe(":8081", router.GetHandleRequest())
}
