package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Scheme string

const (
	HTTP  Scheme = "http"
	HTTPS Scheme = "https"
)

type Config struct {
	Server struct {
		IP     string
		Port   int
		Scheme Scheme
	}
	ServerHttps struct {
		CertFilePath string
		KeyFilePath  string
	}
	Mode          string
	RemoteService struct {
		Token                       string
		CardsServiceAddress         string
		ReadonlyCardsServiceAddress string
	}
	DatabseConnection string
	Validators        []struct {
		AppID         string
		PublicKeyPath string
		PublicKey     string
	}
	LogFile string
}

func ReadConfiguration() {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open config file by '%v' path\n", configPath)
		os.Exit(2)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read config file by '%v' path\n", configPath)
		os.Exit(2)
	}
}
