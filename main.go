package main

import (
	"flag"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/app"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Custom config")
}

func main() {
	flag.Parse()
	app.Init(configPath)
	panic(app.Run())
}
