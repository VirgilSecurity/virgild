package app

import (
	"encoding/hex"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/auth"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/controllers"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols/http"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/signer"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/validators"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/enums"
	"io/ioutil"
	"log"
	"os"
)

var config settings
var server protocols.Server
var logger *log.Logger

func Init(configPath string) {
	readConfiguration(configPath)
	logger = makeLogger()

	storage := makeStorage()
	serviceSigner := makeServiceSigner()

	controller := &controllers.Controller{
		Storage: storage,
		Signer:  serviceSigner,
	}
	authHandler := &auth.AuthHander{
		Token: config.AuthService.Token,
	}

	server = http.MakeServer(config.Server.Host, config.ServerHttps.CertFilePath, config.ServerHttps.KeyFilePath, controller, authHandler)
}

func Run() error {
	return server.Serve()
}

func makeLogger() *log.Logger {
	if config.LogFile != "" {
		return log.New(os.Stderr, "", log.LstdFlags)
	} else {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
}

func makeSignValidator() *validators.SignValidator {
	keys := make(map[string][]byte, 0)
	for _, v := range config.Validators {
		if v.PublicKey != "" {
			keys[v.AppID] = []byte(v.PublicKey)
		} else if v.PublicKeyPath != "" {
			b, err := ioutil.ReadFile(v.PublicKeyPath)
			if err != nil {
				makeLogger().Println("Cannot read file by", v.PublicKeyPath, "path")
				continue
			}
			keys[v.AppID] = b
		}
	}
	return validators.MakeSignValidator(keys)
}

func makeServiceSigner() *signer.ServiceSigner {
	crypto := virgil.Crypto()
	pk, err := ioutil.ReadFile(config.ServiceSigner.PrivateKeyPath)
	if err != nil {
		logger.Println("Cannot read private key for service")
		os.Exit(2)
		return nil
	}
	private, err := crypto.ImportPrivateKey(pk, config.ServiceSigner.Password)
	if err != nil {
		logger.Println("Cannot import private key")
		os.Exit(2)
		return nil
	}
	pubKey, err := private.ExtractPublicKey()
	if err != nil {
		logger.Println("Cannot extract public key from private")
		os.Exit(2)
		return nil
	}
	vr, err := virgil.NewCreateCardRequest("application", config.ServiceSigner.Identity, pubKey, enums.CardScope.Global, nil)
	if err != nil {
		logger.Println("Cannot create card request for service signer")
		os.Exit(2)
		return nil
	}
	requestSigner := &virgil.RequestSigner{}
	requestSigner.SelfSign(vr, private)

	snapshot, _ := vr.GetSnapshot()
	id := hex.EncodeToString(crypto.CalculateFingerprint(snapshot))
	s := makeLocalStorage()

	r, errResp := s.GetCard(id)
	if errResp != nil {
		logger.Println("Cannot add self card into local storage")
		os.Exit(2)
		return nil
	}
	if r == nil {
		s.CreateCard(&models.CardResponse{
			Snapshot: snapshot,
			Meta: models.ResponseMeta{
				Signatures: vr.Signatures,
			},
		})
	}

	return &signer.ServiceSigner{
		ID:         id,
		PrivateKey: private,
	}
}
