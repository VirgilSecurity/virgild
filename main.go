package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/http"

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
