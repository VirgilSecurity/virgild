package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type settings struct {
	Server struct {
		Host string
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
	LogFile     string
	AuthService struct {
		Token string
	}
	ServiceSigner struct {
		Identity       string
		PrivateKeyPath string
		Password       string
	}
}

func readConfiguration(configPath string) {
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
