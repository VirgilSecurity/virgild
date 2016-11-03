package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/http"
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

	router := http.MakeRouter(&controllers.Controller{
		Storage: storage,
		// Validator: MakeSignValidator(),
	}, MakeLogger())

	fasthttp.ListenAndServe(":8081", router.GetHandleRequest())
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
