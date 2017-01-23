package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-xorm/xorm"

	"gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/transport/virgilhttp"
	"gopkg.in/virgil.v4/virgilcrypto"
)

type SignerConf struct {
	CardID     string `json:"card_id"`
	PrivateKey string `json:"private_key"`
}

type VRAConf struct {
	CardID    string `json:"card_id"`
	PublicKey string `json:"public_key"`
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
	Token    string          `json:"token,omitempty"`
	Cache    int             `json:"cache,omitempty"`
	Services ServicesConf    `json:"services,omitempty"`
	Trust    RemoteTrustConf `json:"trust,omitempty"`
}

type Config struct {
	DB      string      `json:"db"`
	Signer  SignerConf  `json:"signer"`
	VRAConf *VRAConf    `json:"vra,omitempty"`
	Remote  *RemoteConf `json:"remote,omitempty"`
}

type SignerConfig struct {
	CardID     string
	PrivateKey virgilcrypto.PrivateKey
}

type VRAConfig struct {
	CardID    string
	PublicKey virgilcrypto.PublicKey
}

type RemoteConfig struct {
	Token  string
	Cache  time.Duration
	Client *virgil.Client
}
type AppConfig struct {
	Orm        *xorm.Engine
	Logger     *log.Logger
	Crypto     virgilcrypto.Crypto
	Signer     SignerConfig
	Remote     *RemoteConfig
	VRA        *VRAConfig
	ConfigPath string
	Config     Config
}

var (
	app AppConfig
)

func loadAppConfig(configPath string) *AppConfig {
	// virgilcrypto.DefaultCrypto = &crypto.NativeCrypto{}
	app = AppConfig{
		Logger:     log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		Crypto:     virgil.Crypto(),
		ConfigPath: configPath,
	}
	oldConf := loadFromFile()
	app.Config = oldConf
	setupConfig()
	setupDB()
	setupSigner()
	setupRemote()
	setupVRA()

	if oldConf != app.Config {
		saveConfig(app.Config, app.ConfigPath)
	}

	return &app
}

func loadFromFile() Config {
	config := Config{}
	if _, err := os.Stat(app.ConfigPath); err == nil {
		d, err := ioutil.ReadFile(app.ConfigPath)
		if err != nil {
			app.Logger.Fatalf("Cannot read configuration: %v", err)
		}
		err = json.Unmarshal(d, &config)
		if err != nil {
			app.Logger.Fatalf("Cannot load configuration: %v", err)
		}
	}
	return config
}

func setupConfig() {
	if len(app.Config.DB) == 0 {
		app.Config.DB = "sqlite3:./virgild.db"
	}

	if len(app.Config.Signer.PrivateKey) == 0 {
		app.Config.Signer.PrivateKey = "./private.key"
	}
	if _, err := os.Stat(app.Config.Signer.PrivateKey); os.IsNotExist(err) {
		kp, err := app.Crypto.GenerateKeypair()
		if err != nil {
			app.Logger.Fatalf("Cannot generate private key: %v", err)
		}
		d, err := kp.PrivateKey().Encode([]byte(""))
		if err != nil {
			app.Logger.Fatalf("Cannot generate private key: %v", err)
		}
		err = ioutil.WriteFile(app.Config.Signer.PrivateKey, d, 400)
		if err != nil {
			app.Logger.Fatalf("Cannot save private key: %v", err)
		}
	}
}

func setupDB() {
	db := app.Config.DB
	index := strings.Index(db, ":")
	if index == -1 {
		app.Logger.Fatalf("Database connection string incorrect. Expected {provider}:{connection_string} actual: %v", db)
	}

	driver, conn := db[:index], db[index+1:]
	orm, err := xorm.NewEngine(driver, conn)
	if err != nil {
		app.Logger.Fatalf("Cannot connect to db (Provider: %v Connection: %v): %v", driver, conn, err)
	}
	if err = Sync(orm); err != nil {
		app.Logger.Fatalf("Cannot sync tabls: %v", err)
	}
	app.Orm = orm
}

func setupSigner() {
	d, err := ioutil.ReadFile(app.Config.Signer.PrivateKey)
	if err != nil {
		app.Logger.Fatalf("Cannot load private key : %+v", err)
	}
	privateKey, err := app.Crypto.ImportPrivateKey(d, "")
	if err != nil {
		app.Logger.Fatalf("Unsupporetd format of private key: %v", err)
	}

	if len(app.Config.Signer.CardID) == 0 {
		app.Config.Signer.CardID = createCard(privateKey)
	}
	app.Signer.CardID = app.Config.Signer.CardID
	app.Signer.PrivateKey = privateKey
}

func createCard(key virgilcrypto.PrivateKey) string {
	pub, err := key.ExtractPublicKey()
	if err != nil {
		app.Logger.Fatalf("Cannot extract public key: %v", err)
	}

	info := virgil.CardModel{
		Identity:     "virgild",
		IdentityType: "card service",
		Scope:        virgil.CardScope.Application,
	}
	req, err := virgil.NewCreateCardRequest(info.Identity, info.IdentityType, pub, virgil.CardParams{
		Scope: info.Scope,
	})
	if err != nil {
		app.Logger.Fatalf("Cannot create card for virgild: %+v", err)
	}
	id := hex.EncodeToString(app.Crypto.CalculateFingerprint(req.Snapshot))
	epub, err := pub.Encode()
	if err != nil {
		app.Logger.Fatalf("Cannot encode public key: %v", err)
	}
	fmt.Println("Public Key:", base64.StdEncoding.EncodeToString(epub))
	fmt.Println("ID:", id)

	h := DefaultModeCardHandler{
		Repo: &ImpSqlCardRepository{
			Orm: app.Orm,
		},
		Fingerprint: &ImpFingerprint{
			Crypto: app.Crypto,
		},
		Signer: &ImpRequestSigner{
			CardId:     id,
			PrivateKey: key,
			Crypto:     app.Crypto,
		},
		Validator: &RequestValidator{
			CreateCardValidators: make([]func(req *CreateCardRequest) (bool, error), 0),
		},
	}

	_, err = h.Create(&CreateCardRequest{
		Info:    info,
		Request: *req,
	})
	if err != nil {
		app.Logger.Fatalf("Cannot store virgild card: %+v", err)
	}
	return id
}

func setupRemote() {
	if app.Config.Remote == nil {
		return
	}

	if len(app.Config.Remote.Token) == 0 {
		app.Logger.Fatalf("Remote connection requared set token")
	}
	if len(app.Config.Remote.Services.Cards) == 0 {
		app.Config.Remote.Services.Cards = "https://cards.virgilsecurity.com"
	}
	if len(app.Config.Remote.Services.CardsRO) == 0 {
		app.Config.Remote.Services.CardsRO = "https://cards-ro.virgilsecurity.com"
	}
	if len(app.Config.Remote.Services.Identity) == 0 {
		app.Config.Remote.Services.Identity = "https://identity.virgilsecurity.com"
	}
	if app.Config.Remote.Cache == 0 {
		app.Config.Remote.Cache = 3600 // 1 hour
	}
	if len(app.Config.Remote.Trust.PublicKey) == 0 {
		app.Config.Remote.Trust.PublicKey = "./trust.pub"
	}
	if _, err := os.Stat(app.Config.Remote.Trust.PublicKey); err != nil {
		err = ioutil.WriteFile(app.Config.Remote.Trust.PublicKey, []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEA8jJqWY5hm4tvmnM6QXFdFCErRCnoYdhVNjFggffSCoc=
-----END PUBLIC KEY-----`), 700)
		if err != nil {
			app.Logger.Fatalf("Cannot save trust service public key: %+v", err)
		}
		app.Config.Remote.Trust.CardID = "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853"
	}

	customValidator := virgil.NewCardsValidator()
	d, err := ioutil.ReadFile(app.Config.Remote.Trust.PublicKey)
	if err != nil {
		app.Logger.Fatalf("Cannot load public key of trust service: %+v", err)
	}
	cardsServicePublic, err := app.Crypto.ImportPublicKey(d)
	if err != nil {
		app.Logger.Fatalf("Cannot import public key: %v", err)
	}
	customValidator.AddVerifier(app.Config.Remote.Trust.CardID, cardsServicePublic)
	//""

	vclient, err := virgil.NewClient(app.Config.Remote.Token,
		virgil.ClientTransport(
			virgilhttp.NewTransportClient(
				app.Config.Remote.Services.Cards,
				app.Config.Remote.Services.CardsRO,
				app.Config.Remote.Services.Identity,
				app.Config.Remote.Services.Identity)), // VRA
		virgil.ClientCardsValidator(customValidator))

	app.Remote = &RemoteConfig{
		Token:  app.Config.Remote.Token,
		Client: vclient,
		Cache:  time.Duration(app.Config.Remote.Cache) * time.Second,
	}
}

func saveConfig(config Config, configPath string) {
	d, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		app.Logger.Fatalf("Cannot serialaze configuration: %v", err)
	}
	if err = ioutil.WriteFile(configPath, d, 400); err != nil {
		app.Logger.Fatalf("Cannot save configuration: %v", err)
	}
}

func setupVRA() {
	if app.Config.VRAConf == nil {
		return
	}
	if len(app.Config.VRAConf.CardID) == 0 {
		app.Logger.Fatalf("VRA card id is not set")
	}

	d, err := ioutil.ReadFile(app.Config.VRAConf.PublicKey)
	if err != nil {
		app.Logger.Fatalf("Cannot load public key of VRA service: %+v", err)
	}
	vraPubKey, err := app.Crypto.ImportPublicKey(d)
	if err != nil {
		app.Logger.Fatalf("Cannot import VRA public key: %v", err)
	}

	app.VRA = &VRAConfig{
		CardID:    app.Config.VRAConf.CardID,
		PublicKey: vraPubKey,
	}
}
