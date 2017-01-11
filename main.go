package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/valyala/fasthttp"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/transport/virgilhttp"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/virgilcrypto"
)

var (
	configPath string
	logger     *log.Logger
	config     *Config
	orm        *xorm.Engine
	crypto     virgilcrypto.Crypto
)

func init() {
	flag.StringVar(&configPath, "config", "./virgild.conf", "Configuration file")
}

func main() {
	flag.Parse()
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	loadConfig()
	setupDB()

	crypto = virgil.Crypto()

	d, err := ioutil.ReadFile(config.Signer.PrivateKey)
	if err != nil {
		logger.Fatalf("Cannot load private key : %+v", err)
	}
	privateKey, err := crypto.ImportPrivateKey(d, "")
	if err != nil {
		logger.Fatalf("Unsupporetd format of private key: %v", err)
	}

	if len(config.Signer.CardID) == 0 {
		config.Signer.CardID = createCard(privateKey)
	}

	validator := &RequestValidator{
		SearchValidators: []func(criteria *virgil.Criteria) (bool, error){
			ScopeMustGlobalOrApplication,
			SearchIdentitiesNotEmpty,
		},
		CreateCardValidators: []func(req *CreateCardRequest) (bool, error){
			CardIdentityIsEmpty,
			CardPublicKeyLengthInvalid,
			CreateCardRequestSignsEmpty,
			CardDataEntries,
			CardDataValueExceed256,
			CardInfoValueExceed256,
			GlobalCardIdentityTypeMustBeEmail,
		},
		RevokeCardValidators: []func(req *RevokeCardRequest) (bool, error){
			//		RevokeCardRequestSignsEmpty,
			RevocationReasonIsInvalide,
		},
	}

	var router Router

	if config.Remote != nil {
		setupRemoteConf()

		customValidator := virgil.NewCardsValidator()
		d, err = ioutil.ReadFile(config.Remote.Trust.PublicKey)
		if err != nil {
			logger.Fatalf("Cannot load public key of trust service: %+v", err)
		}
		cardsServicePublic, err := crypto.ImportPublicKey(d)
		if err != nil {
			logger.Fatalf("Cannot import public key: %v", err)
		}
		customValidator.AddVerifier(config.Remote.Trust.CardID, cardsServicePublic)
		//""

		vclient, err := virgil.NewClient(config.Remote.Token,
			virgil.ClientTransport(
				virgilhttp.NewTransportClient(
					config.Remote.Services.Cards,
					config.Remote.Services.CardsRO,
					config.Remote.Services.Identity)),
			virgil.ClientCardsValidator(customValidator))

		router = Router{
			Card: &CardController{
				Logger: logger,
				Card: &AppModeCardHandler{
					Repo: &ImpSqlCardRepository{
						Cache: time.Duration(config.Remote.CacheDuration) * time.Second,
						Orm:   orm,
					},
					Signer: &ImpRequestSigner{
						CardId:     config.Signer.CardID,
						PrivateKey: privateKey,
						Crypto:     crypto,
					},
					Validator: validator,
					Remote:    vclient,
				},
			},
		}

	} else {
		router = Router{
			Card: &CardController{
				Logger: logger,
				Card: &DefaultModeCardHandler{
					Repo: &ImpSqlCardRepository{
						Orm: orm,
					},
					Fingerprint: &ImpFingerprint{
						Crypto: crypto,
					},
					Signer: &ImpRequestSigner{
						CardId:     config.Signer.CardID,
						PrivateKey: privateKey,
						Crypto:     crypto,
					},
					Validator: validator,
				},
			},
		}
	}

	saveConfig()
	if err != nil {
		logger.Fatalf("Cannot save configuration: %+v", err)
	}

	err = fasthttp.ListenAndServe(":3001", router.Handler())
	if err != nil {
		logger.Fatalln(err)
	}
}

func setupDB() {
	var err error
	db := config.DB
	index := strings.Index(db, ":")
	if index == -1 {
		logger.Fatalf("Database connection string incorrect. Expected {provider}:{connection_string} actual: %v", config.DB)
	}

	driver, conn := db[:index], db[index+1:]
	orm, err = xorm.NewEngine(driver, conn)
	if err != nil {
		logger.Fatalf("Cannot connect to db (Provider: %v Connection: %v): %v", driver, conn, err)
	}
	if err = Sync(orm); err != nil {
		logger.Fatalf("Cannot sync tabls: %v", err)
	}
}

func createCard(key virgilcrypto.PrivateKey) string {
	pub, err := key.ExtractPublicKey()
	if err != nil {
		logger.Fatalf("Cannot extract public key: %v", err)
	}

	info := virgil.CreateCardRequest{
		Identity:     "virgild",
		IdentityType: "card service",
		Scope:        virgil.CardScope.Application,
	}
	req, err := virgil.NewCreateCardRequest(info.Identity, info.IdentityType, pub, virgil.CardInfo{
		Scope: info.Scope,
	})
	if err != nil {
		logger.Fatalf("Cannot create card for virgild: %+v", err)
	}
	id := hex.EncodeToString(crypto.CalculateFingerprint(req.Snapshot))
	epub, err := pub.Encode()
	if err != nil {
		logger.Fatalf("Cannot encode public key: %v", err)
	}
	fmt.Println("Public Key:", base64.StdEncoding.EncodeToString(epub))
	fmt.Println("ID:", id)

	h := DefaultModeCardHandler{
		Repo: &ImpSqlCardRepository{
			Orm: orm,
		},
		Fingerprint: &ImpFingerprint{
			Crypto: crypto,
		},
		Signer: &ImpRequestSigner{
			CardId:     id,
			PrivateKey: key,
			Crypto:     crypto,
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
		logger.Fatalf("Cannot store virgild card: %+v", err)
	}
	return id
}
