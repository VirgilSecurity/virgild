package config

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
	Signer *Signer
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
	DB           *xorm.Engine
	Logger       *log.Logger
	config       Config
	ConfigUpdate *Updater
	ConfigPath   string
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
	app.Common.config = conf
	app.Common.ConfigPath = file

	app.Common.DB, err = initDB(conf.DB)
	if err != nil {
		panic(err)
	}
	app.Common.Logger, err = initLogger(conf.LogFile)
	if err != nil {
		panic(err)
	}
	app.Cards.VRA, err = initVRA(conf.Cards.VRA)
	if err != nil {
		panic(err)
	}
	app.Cards.Signer, err = initSigner(&conf.Cards.Signer)
	if err != nil {
		panic(err)
	}
	app.Cards.Remote, err = initRemote(conf.Cards.Remote)
	if err != nil {
		panic(err)
	}
	app.Site.VirgilD, err = setSiteVirgilD(app.Cards.Signer)
	if err != nil {
		panic(err)
	}
	app.Site.Admin.Login = conf.Admin.Login
	app.Site.Admin.Password = conf.Admin.Password

	if app.Cards.Signer.Card != nil { // first start
		app.Common.config = conf
		saveConfigToFole(conf, file)
	}
	app.Common.ConfigUpdate = &Updater{
		app: app,
	}

	return app
}

func initDB(db string) (*xorm.Engine, error) {

	i := strings.Index(db, ":")
	if i < 0 {
		return nil, fmt.Errorf("Cannot pars database connectin ({provider:connection})")
	}
	d, c := db[:i], db[i+1:]
	e, err := xorm.NewEngine(d, c)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to (driver: %v name: %v) database: %v", d, c, err)
	}
	return e, nil
}

func initLogger(logFile string) (*log.Logger, error) {
	var w io.Writer
	if logFile == "console" {
		w = os.Stderr
	} else {
		f, err := os.OpenFile(logFile, os.O_APPEND, 700)
		if err != nil {
			return nil, fmt.Errorf("Cannot open file config (%v): %v", logFile, err)
		}
		w = f
	}
	return log.New(w, "", log.LUTC|log.LstdFlags), nil
}

func initVRA(conf *AuthorityConfig) (*Authority, error) {
	if conf == nil || conf.CardID == "" || conf.PublicKey == "" {
		return nil, nil
	}
	pub, err := virgil.Crypto().ImportPublicKey([]byte(conf.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("Cannot import public key of VRA: %v", err)
	}
	return &Authority{
		CardID:    conf.CardID,
		PublicKey: pub,
	}, nil
}

func initSigner(conf *SignerConfig) (*Signer, error) {
	priv, err := virgil.Crypto().ImportPrivateKey([]byte(conf.PrivateKey), conf.PrivateKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("Cannot load private key for VirgilD: %v", err)
	}
	signer := &Signer{
		PrivateKey: priv,
	}
	if conf.CardID == "" {
		signer.Card, err = createVirgilCard(priv)
		if err != nil {
			return nil, err
		}
		conf.CardID = signer.Card.ID
	}

	signer.CardID = conf.CardID
	return signer, nil
}

func createVirgilCard(key virgilcrypto.PrivateKey) (*virgil.Card, error) {
	pub, err := key.ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf("Cannot extract public key: %v", err)
	}
	req, err := virgil.NewCreateCardRequest("VirgilD", "Service", pub, virgil.CardParams{})
	if err != nil {
		return nil, fmt.Errorf("Cannot create card request: %v", err)
	}
	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, key)
	if err != nil {
		return nil, fmt.Errorf("Cannot self sign VirgilD card: %v", err)
	}
	id := hex.EncodeToString(virgil.Crypto().CalculateFingerprint(req.Snapshot))
	return &virgil.Card{
		ID:           id,
		Identity:     "VirgilD",
		IdentityType: "Service",
		Scope:        virgil.CardScope.Application,
		Snapshot:     req.Snapshot,
		Signatures:   req.Meta.Signatures,
	}, nil
}

func initRemote(conf *RemoteConfig) (*Remote, error) {
	if conf == nil {
		return nil, nil
	}

	customValidator := virgil.NewCardsValidator()

	cardsServicePublic, err := virgil.Crypto().ImportPublicKey([]byte(conf.Authority.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("Cannot load public key of Authority Service: %+v", err)
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
		return nil, fmt.Errorf("Cannot init Virgil Client: %+v", err)
	}
	return &Remote{
		VClient: vclient,
		Cache:   time.Duration(conf.Cache) * time.Second,
	}, nil
}

func setSiteVirgilD(signer *Signer) (VirgilDCard, error) {
	pub, err := signer.PrivateKey.ExtractPublicKey()
	if err != nil {
		return VirgilDCard{}, fmt.Errorf("Cannot extract public key from VirgilD's private key: %+v", err)
	}
	b, err := pub.Encode()
	if err != nil {
		return VirgilDCard{}, fmt.Errorf("Cannot encpde VirgilD's public key: %+v", err)
	}
	return VirgilDCard{
		CardID:    signer.CardID,
		PublicKey: base64.StdEncoding.EncodeToString(b),
	}, nil
}
