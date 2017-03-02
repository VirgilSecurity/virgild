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

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tochka/flag"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/transport/virgilhttp"
	"gopkg.in/virgil.v4/virgilcrypto"
	crypto "gopkg.in/virgilsecurity/virgil-crypto-go.v4"
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

type CardMode string

const (
	CardModeCache CardMode = "cache"
	CardModeLocal CardMode = "local"
	CardModeSync  CardMode = "sync"
)

type Cards struct {
	Mode   CardMode
	Signer *Signer
	VRA    *Authority
	Remote Remote
}

type SiteAdmin struct {
	Login    string
	Password string
}

type VirgilDCard struct {
	CardID    string `json:"card_id"`
	PublicKey string `json:"public_key"`
}
type Site struct {
	Enabled bool
	Admin   SiteAdmin
	VirgilD VirgilDCard
}

type Common struct {
	DB           *xorm.Engine
	Logger       *log.Logger
	Config       Config
	ConfigUpdate *Updater
	ConfigPath   string
	Address      string
}

type AuthMode string

const (
	AuthModeLocal    AuthMode = "local"
	AuthModeExternal AuthMode = "external"
	AuthModeNo       AuthMode = "no"
)

type AuthParams struct {
	Host string
}

type Auth struct {
	Mode      AuthMode
	TokenType string
	Params    AuthParams
}

type App struct {
	Site   Site
	Common Common
	Cards  Cards
	Auth   Auth
}

func Init() *App {
	var err error
	app := new(App)
	virgilcrypto.DefaultCrypto = &crypto.NativeCrypto{}

	flag.Parse()

	conf := defaultConfig
	app.Common.Config = conf
	if configPath == "" {
		app.Common.ConfigPath = "virgild.conf"
	} else {
		app.Common.ConfigPath = configPath
	}

	app.Common.Address = conf.Address

	app.Common.DB, err = initDB(conf.DB)
	if err != nil {
		panic(err)
	}
	app.Common.Logger, err = initLogger(conf.LogFile)
	if err != nil {
		panic(err)
	}
	app.Cards, err = initCards(&conf.Cards)
	if err != nil {
		panic(err)
	}

	app.Site.VirgilD, err = setSiteVirgilD(app.Cards.Signer)
	if err != nil {
		panic(err)
	}
	app.Site.Admin.Login = conf.Admin.Login
	app.Site.Admin.Password = conf.Admin.Password
	app.Site.Enabled = conf.Admin.Enabled

	app.Auth, err = initAtuh(conf.Auth)
	if err != nil {
		panic(err)
	}
	if app.Common.Config != conf { // has changes
		app.Common.Config = conf
		saveConfigToFole(conf, app.Common.ConfigPath)
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

func initCards(conf *CardsConfig) (cards Cards, err error) {
	switch CardMode(conf.Mode) {
	case CardModeCache, CardMode(""):
		cards.Mode = CardModeCache
	case CardModeLocal, CardModeSync:
		cards.Mode = CardMode(conf.Mode)
	default:
		err = fmt.Errorf("Unsupported cards mode (%v)", conf.Mode)
		return
	}
	cards.VRA, err = initVRA(conf.VRA)
	if err != nil {
		return
	}
	if cards.Mode != CardModeCache {
		cards.Signer, err = initSigner(&conf.Signer)
		if err != nil {
			return
		}
	}
	cards.Remote, err = initRemote(conf.Remote)
	if err != nil {
		return
	}
	return
}

func initVRA(conf AuthorityConfig) (*Authority, error) {
	if conf.CardID == "" || conf.PublicKey == "" {
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
	if conf.PrivateKey == "" {
		pk, err := virgil.Crypto().GenerateKeypair()
		if err != nil {
			log.Fatalf("Cannot generate keypair for VirgilD: %v", err)
		}
		p, err := pk.PrivateKey().Encode([]byte(""))
		if err != nil {
			log.Fatalf("Cannot encode private key for VirgilD: %v", err)
		}

		conf.PrivateKey = base64.StdEncoding.EncodeToString(p)
	}

	b, err := base64.StdEncoding.DecodeString(conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode authority public key from base64: %+v", err)
	}
	priv, err := virgil.Crypto().ImportPrivateKey(b, conf.PrivateKeyPassword)
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

func initRemote(conf RemoteConfig) (Remote, error) {
	customValidator := virgil.NewCardsValidator()

	b, err := base64.StdEncoding.DecodeString(conf.Authority.PublicKey)
	if err != nil {
		return Remote{}, fmt.Errorf("Cannot decode authority public key from base64: %+v", err)
	}
	cardsServicePublic, err := virgil.Crypto().ImportPublicKey(b)
	if err != nil {
		return Remote{}, fmt.Errorf("Cannot load public key of Authority Service: %+v", err)
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
		return Remote{}, fmt.Errorf("Cannot init Virgil Client: %+v", err)
	}
	return Remote{
		VClient: vclient,
		Cache:   time.Duration(conf.Cache) * time.Second,
	}, nil
}

func setSiteVirgilD(signer *Signer) (VirgilDCard, error) {
	if signer == nil {
		return VirgilDCard{}, nil
	}

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

func initAtuh(conf AuthConfig) (Auth, error) {
	t := conf.TokenType
	if t == "" {
		t = "VIRGIL"
	}
	switch AuthMode(conf.Mode) {
	case AuthModeLocal:
		return Auth{
			Mode:      AuthModeLocal,
			TokenType: t,
		}, nil
	case AuthModeNo, AuthMode(""):
		return Auth{
			Mode:      AuthModeNo,
			TokenType: t,
		}, nil
	case AuthModeExternal:
		if conf.Params.Host == "" {
			return Auth{}, fmt.Errorf("Auth config invalid. For external mode auth must be set the host of Auth service")
		}
		return Auth{
			Mode:      AuthModeExternal,
			TokenType: t,
			Params: AuthParams{
				Host: conf.Params.Host,
			},
		}, nil
	default:
		return Auth{}, fmt.Errorf("Undefined auth mode (%v). Supported [no, local, external]", conf.Mode)
	}
}
