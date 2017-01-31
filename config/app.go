package config

import (
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/transport/virgilhttp"
	"gopkg.in/virgil.v4/virgilcrypto"
)

type Signer struct {
	CardID     string
	Card       *virgil.Card
	PrivateKey virgilcrypto.PrivateKey
}

type Authority struct {
	CardID    string
	PublicKey virgilcrypto.PublicKey
}

type Remote struct {
	Cache   time.Duration
	VClient *virgil.Client
}

type Cards struct {
	Signer Signer
	VRA    *Authority
	Remote *Remote
}

type SiteAdmin struct {
	Login    string
	Password string
}

type VirgilDCard struct {
	CardID    string
	PublicKey string
}
type Site struct {
	Admin   SiteAdmin
	VirgilD VirgilDCard
}

type Common struct {
	DB         *xorm.Engine
	Logger     *log.Logger
	oldConfig  Config
	Config     Config
	ConfigPath string
}

type App struct {
	Site   Site
	Common Common
	Cards  Cards
}

func Init(file string) *App {
	app := new(App)

	conf, err := loadConfigFromFile(file)
	if err != nil {
		log.Fatalf("Cannot load configuration file: %v", err)
	}

	conf = initDefault(conf)
	app.Common.oldConfig = conf
	app.Common.ConfigPath = file

	app.Common.DB = initDB(conf.DB)
	app.Common.Logger = initLogger(conf.LogFile)
	app.Cards.VRA = initVRA(conf.Cards.VRA)
	app.Cards.Signer = initSigner(&conf.Cards.Signer)
	app.Cards.Remote = initRemote(conf.Cards.Remote)

	app.Site.VirgilD = setSiteVirgilD(app.Cards.Signer)
	app.Site.Admin.Login = conf.Admin.Login
	app.Site.Admin.Password = conf.Admin.Password

	if app.Cards.Signer.Card != nil { // first start
		app.Common.Config = conf
		saveConfigToFole(conf, file)
	}
	return app
}

func initDB(db string) *xorm.Engine {

	i := strings.Index(db, ":")
	if i < 0 {
		log.Fatalf("Cannot pars database connectin ({provider:connection})")
	}
	d, c := db[:i], db[i+1:]
	e, err := xorm.NewEngine(d, c)
	if err != nil {
		log.Fatalf("Cannot connect to (driver: %v name: %v) database: %v", d, c, err)
	}
	return e
}

func initLogger(logFile string) *log.Logger {
	var w io.Writer
	if logFile == "console" {
		w = os.Stderr
	} else {
		f, err := os.OpenFile(logFile, os.O_APPEND, 700)
		if err != nil {
			log.Fatalf("Cannot open file config (%v): %v", logFile, err)
		}
		w = f
	}
	return log.New(w, "", log.LUTC|log.LstdFlags)
}

func initVRA(conf *AuthorityConfig) *Authority {
	if conf == nil || conf.CardID == "" || conf.PublicKey == "" {
		return nil
	}
	pub, err := virgil.Crypto().ImportPublicKey([]byte(conf.PublicKey))
	if err != nil {
		log.Fatalf("Cannot import public key of VRA: %v", err)
	}
	return &Authority{
		CardID:    conf.CardID,
		PublicKey: pub,
	}
}

func initSigner(conf *SignerConfig) Signer {
	priv, err := virgil.Crypto().ImportPrivateKey([]byte(conf.PrivateKey), conf.PrivateKeyPassword)
	if err != nil {
		log.Fatalf("Cannot load private key for VirgilD: %v", err)
	}
	signer := Signer{
		PrivateKey: priv,
	}
	if conf.CardID == "" {
		signer.Card = createVirgilCard(priv)
		conf.CardID = signer.Card.ID
	}

	signer.CardID = conf.CardID
	return signer
}

func createVirgilCard(key virgilcrypto.PrivateKey) *virgil.Card {
	pub, err := key.ExtractPublicKey()
	if err != nil {
		log.Fatalf("Cannot extract public key: %v", err)
	}
	req, err := virgil.NewCreateCardRequest("VirgilD", "Service", pub, virgil.CardParams{})
	if err != nil {
		log.Fatalf("Cannot create card request: %v", err)
	}
	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, key)
	if err != nil {
		log.Fatalf("Cannot self sign VirgilD card: %v", err)
	}
	id := hex.EncodeToString(virgil.Crypto().CalculateFingerprint(req.Snapshot))
	return &virgil.Card{
		ID:           id,
		Identity:     "VirgilD",
		IdentityType: "Service",
		Scope:        virgil.CardScope.Application,
		Snapshot:     req.Snapshot,
		Signatures:   req.Meta.Signatures,
	}
}

func initRemote(conf *RemoteConfig) *Remote {
	if conf == nil {
		return nil
	}

	customValidator := virgil.NewCardsValidator()

	cardsServicePublic, err := virgil.Crypto().ImportPublicKey([]byte(conf.Authority.PublicKey))
	if err != nil {
		log.Fatalf("Cannot load public key of Authority Service: %+v", err)
	}

	customValidator.AddVerifier(conf.Authority.CardID, cardsServicePublic)
	vclient, err := virgil.NewClient(conf.Token,
		virgil.ClientTransport(
			virgilhttp.NewTransportClient(
				conf.Services.Cards,
				conf.Services.CardsRO,
				conf.Services.Identity,
				conf.Services.VRA)),
		virgil.ClientCardsValidator(customValidator))
	if err != nil {
		log.Fatalf("Cannot init Virgil Client: %+v", err)
	}
	return &Remote{
		VClient: vclient,
		Cache:   time.Duration(conf.Cache) * time.Second,
	}
}

func setSiteVirgilD(signer Signer) VirgilDCard {
	pub, err := signer.PrivateKey.ExtractPublicKey()
	if err != nil {
		log.Fatalf("Cannot extract public key from VirgilD's private key: %+v", err)
	}
	b, err := pub.Encode()
	if err != nil {
		log.Fatalf("Cannot encpde VirgilD's public key: %+v", err)
	}
	return VirgilDCard{
		CardID:    signer.CardID,
		PublicKey: base64.StdEncoding.EncodeToString(b),
	}
}
