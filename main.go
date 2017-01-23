package main

import (
	"flag"

	"github.com/valyala/fasthttp"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./virgild.conf", "Configuration file")
}

func main() {
	flag.Parse()
	appConfig := loadAppConfig(configPath)

	router := Router{
		Card: &CardController{
			MakeResponse: MakeResponse(appConfig.Logger),
			Card:         getCardHandler(appConfig),
		},
	}

	err := fasthttp.ListenAndServe(":3001", router.Handler())
	if err != nil {
		app.Logger.Fatalln(err)
	}
}

func getCardHandler(appConfig *AppConfig) CardHandler {
	if appConfig.Remote != nil {
		return &AppModeCardHandler{
			Repo: &ImpSqlCardRepository{
				Cache: appConfig.Remote.Cache,
				Orm:   appConfig.Orm,
			},
			Signer: &ImpRequestSigner{
				CardId:     appConfig.Signer.CardID,
				PrivateKey: appConfig.Signer.PrivateKey,
				Crypto:     appConfig.Crypto,
			},
			Validator: MakeRequestValidation(appConfig),
			Remote:    appConfig.Remote.Client,
		}
	}
	return &DefaultModeCardHandler{
		Repo: &ImpSqlCardRepository{
			Orm: appConfig.Orm,
		},
		Fingerprint: &ImpFingerprint{
			Crypto: appConfig.Crypto,
		},
		Signer: &ImpRequestSigner{
			CardId:     appConfig.Signer.CardID,
			PrivateKey: appConfig.Signer.PrivateKey,
			Crypto:     appConfig.Crypto,
		},
		Validator: MakeRequestValidation(appConfig),
	}
}
