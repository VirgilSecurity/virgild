package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type SignerConf struct {
	CardID     string `json:"card_id"`
	PrivateKey string `json:"private_key"`
}

type ServicesConf struct {
	Cards    string `json:"cards"`
	CardsRO  string `json:"cards_ro"`
	Identity string `json:"identity"`
}

type RemoteTrustConf struct {
	CardID    string `json:"card_id,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
}

type RemoteConf struct {
	Token         string          `json:"token,omitempty"`
	CacheDuration int             `json:"cache_duration,omitempty"`
	Services      ServicesConf    `json:"services,omitempty"`
	Trust         RemoteTrustConf `json:"trust,omitempty"`
}

type Config struct {
	DB     string      `json:"db"`
	Signer SignerConf  `json:"signer"`
	Remote *RemoteConf `json:"remote,omitempty"`
}

func loadConfig() {
	config = &Config{
		DB: "sqlite3:./virgild.db",
	}
	if _, err := os.Stat(configPath); err == nil {
		d, err := ioutil.ReadFile(configPath)
		if err != nil {
			logger.Fatalf("Cannot read configuration: %v", err)
		}
		err = json.Unmarshal(d, config)
		if err != nil {
			logger.Fatalf("Cannot load configuration: %v", err)
		}
	}

	if len(config.Signer.PrivateKey) == 0 {
		config.Signer.PrivateKey = "./private.key"
	}
	if _, err := os.Stat(config.Signer.PrivateKey); os.IsNotExist(err) {
		kp, err := virgil.Crypto().GenerateKeypair()
		if err != nil {
			logger.Fatalf("Cannot generate private key: %v", err)
		}
		d, err := kp.PrivateKey().Encode([]byte(""))
		if err != nil {
			logger.Fatalf("Cannot generate private key: %v", err)
		}
		err = ioutil.WriteFile(config.Signer.PrivateKey, d, 400)
		if err != nil {
			logger.Fatalf("Cannot save private key: %v", err)
		}
	}
}

func saveConfig() {
	d, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		logger.Fatalf("Cannot serialaze configuration: %v", err)
	}
	if err = ioutil.WriteFile(configPath, d, 400); err != nil {
		logger.Fatalf("Cannot save configuration: %v", err)
	}
}

func setupRemoteConf() {
	if len(config.Remote.Token) == 0 {
		logger.Fatalf("Remote connection requared set token")
	}
	if len(config.Remote.Services.Cards) == 0 {
		config.Remote.Services.Cards = "https://cards.virgilsecurity.com"
	}
	if len(config.Remote.Services.CardsRO) == 0 {
		config.Remote.Services.CardsRO = "https://cards-ro.virgilsecurity.com"
	}
	if len(config.Remote.Services.Identity) == 0 {
		config.Remote.Services.Identity = "https://identity.virgilsecurity.com"
	}
	if config.Remote.CacheDuration == 0 {
		config.Remote.CacheDuration = 3600 // 1 hour
	}
	if len(config.Remote.Trust.PublicKey) == 0 {
		config.Remote.Trust.PublicKey = "./trust.pub"
	}
	if _, err := os.Stat(config.Remote.Trust.PublicKey); err != nil {
		err = ioutil.WriteFile(config.Remote.Trust.PublicKey, []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEA8jJqWY5hm4tvmnM6QXFdFCErRCnoYdhVNjFggffSCoc=
-----END PUBLIC KEY-----`), 700)
		if err != nil {
			logger.Fatalf("Cannot save trust service public key: %+v", err)
		}
		config.Remote.Trust.CardID = "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853"
	}
}
