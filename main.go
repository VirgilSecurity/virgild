package main

import (
	"flag"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/auth"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols/http"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/validators"
	"io/ioutil"

	"log"
	"os"
)

var (
	config     Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Custom config")
}

func main() {
	flag.Parse()

	ReadConfiguration()
	storage := MakeStorage()

	controller := &controllers.Controller{
		Storage: storage,
	}
	authHandler := &auth.AuthHander{
		Token: config.AuthService.Token,
	}
	server := http.MakeServer(config.Server.Host, controller, authHandler)
	panic(server.Serve())
}

func MakeLogger() *log.Logger {
	if config.LogFile != "" {
		return log.New(os.Stderr, "", log.LstdFlags)
	} else {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
}

func MakeSignValidator() *validators.SignValidator {
	keys := make(map[string][]byte, 0)
	for _, v := range config.Validators {
		if v.PublicKey != "" {
			keys[v.AppID] = []byte(v.PublicKey)
		} else if v.PublicKeyPath != "" {
			b, err := ioutil.ReadFile(v.PublicKeyPath)
			if err != nil {
				MakeLogger().Println("Cannot read file by", v.PublicKeyPath, "path")
				continue
			}
			keys[v.AppID] = b
		}
	}
	return validators.MakeSignValidator(keys)
}
