package config

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	virgil "gopkg.in/virgil.v4"
)

type SignerConfig struct {
	CardID             string `json:"card_id"`
	PrivateKey         string `json:"private_key"`
	PrivateKeyPassword string `json:"private_key_password"`
}

type AuthorityConfig struct {
	CardID    string `json:"card_id"`
	PublicKey string `json:"public_key"`
}

type ServicesConfig struct {
	Cards    string `json:"cards"`
	CardsRO  string `json:"cards_ro"`
	Identity string `json:"identity"`
	VRA      string `json:"vra"`
}

type RemoteConfig struct {
	Token     string          `json:"token,omitempty"`
	Cache     int             `json:"cache,omitempty"`
	Services  ServicesConfig  `json:"services,omitempty"`
	Authority AuthorityConfig `json:"authority,omitempty"`
}

type CardsConfig struct {
	Mode   string           `json:"mode"`
	Signer SignerConfig     `json:"signer,omitempty"`
	VRA    *AuthorityConfig `json:"vra,omitempty"`
	Remote RemoteConfig     `json:"remote,omitempty"`
}

type AdminConfig struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthConfig struct {
	Mode      string `json:"mode"`
	TokenType string `json:"token_type"`
	Params    struct {
		Host string `json:"host"`
	} `json:"params"`
}

type Config struct {
	Admin   AdminConfig `json:"admin"`
	DB      string      `json:"db"`
	LogFile string      `json:"log,omitempty"`
	Cards   CardsConfig `json:"cards"`
	Auth    AuthConfig  `json:"auth"`
}

func loadConfigFromFile(file string) (Config, error) {
	var conf Config
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return conf, nil
	}
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(d, &conf)
	return conf, err
}

func saveConfigToFole(config Config, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(config)

	f.Close()
	return err
}

func initDefault(conf Config) Config {
	if conf.DB == "" {
		conf.DB = "sqlite3:virgild.db"
	}
	if conf.LogFile == "" {
		conf.LogFile = "console"
	}
	if conf.Admin.Login == "" {
		conf.Admin.Login = "admin"
	}
	if conf.Admin.Password == "" {
		h := sha256.Sum256([]byte("admin"))
		conf.Admin.Password = hex.EncodeToString(h[:])
	}
	if conf.Cards.Signer.PrivateKey == "" {
		pk, err := virgil.Crypto().GenerateKeypair()
		if err != nil {
			log.Fatalf("Cannot generate keypair for VirgilD: %v", err)
		}
		p, err := pk.PrivateKey().Encode([]byte(""))
		if err != nil {
			log.Fatalf("Cannot encode private key for VirgilD: %v", err)
		}

		conf.Cards.Signer = SignerConfig{
			PrivateKey:         base64.StdEncoding.EncodeToString(p),
			PrivateKeyPassword: "",
			CardID:             "",
		}
	}
	if conf.Cards.Remote.Services.Cards == "" {
		conf.Cards.Remote.Services.Cards = "https://cards.virgilsecurity.com"
	}
	if conf.Cards.Remote.Services.CardsRO == "" {
		conf.Cards.Remote.Services.CardsRO = "https://cards-ro.virgilsecurity.com"
	}
	if conf.Cards.Remote.Services.Identity == "" {
		conf.Cards.Remote.Services.Identity = "https://identity.virgilsecurity.com"
	}
	if conf.Cards.Remote.Services.VRA == "" {
		conf.Cards.Remote.Services.VRA = "https://ra.virgilsecurity.com"
	}
	if conf.Cards.Remote.Authority.CardID == "" || conf.Cards.Remote.Authority.PublicKey == "" {
		conf.Cards.Remote.Authority.CardID = "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853"
		conf.Cards.Remote.Authority.PublicKey = "MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs="
	}
	if conf.Cards.Remote.Cache == 0 {
		conf.Cards.Remote.Cache = 3600
	}
	return conf
}
