package app

import (
	"github.com/virgilsecurity/virgil-apps-cards-cacher/auth"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols/http"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/validators"
	"io/ioutil"
	"log"
	"os"
)

var config settings
var server protocols.Server

func Init(configPath string) {
	readConfiguration(configPath)

	storage := makeStorage()

	controller := &controllers.Controller{
		Storage: storage,
	}
	authHandler := &auth.AuthHander{
		Token: config.AuthService.Token,
	}
	server = http.MakeServer(config.Server.Host, controller, authHandler)
}

func Run() error {
	return server.Serve()
}

func makeLogger() *log.Logger {
	if config.LogFile != "" {
		return log.New(os.Stderr, "", log.LstdFlags)
	} else {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
}

func makeSignValidator() *validators.SignValidator {
	keys := make(map[string][]byte, 0)
	for _, v := range config.Validators {
		if v.PublicKey != "" {
			keys[v.AppID] = []byte(v.PublicKey)
		} else if v.PublicKeyPath != "" {
			b, err := ioutil.ReadFile(v.PublicKeyPath)
			if err != nil {
				makeLogger().Println("Cannot read file by", v.PublicKeyPath, "path")
				continue
			}
			keys[v.AppID] = b
		}
	}
	return validators.MakeSignValidator(keys)
}
