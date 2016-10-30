package remote

import (
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/enums"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/search"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/virgilcrypto"
	"io/ioutil"
)

type RemoteConfig struct {
	CardsServiceAddress         string
	ReadonlyCardsServiceAddress string
	AppID                       string
	Passsword                   string
	AppKey                      []byte
	AppKeyPath                  string
}

func MakeRemoteStorage(token string, conf RemoteConfig) *Remote {
	if conf.CardsServiceAddress == "" {
		conf.CardsServiceAddress = virgil.DefaultClientParams.CardsServiceAddress
	}

	if conf.ReadonlyCardsServiceAddress == "" {
		conf.ReadonlyCardsServiceAddress = virgil.DefaultClientParams.ReadonlyCardsServiceAddress
	}

	virgil.DefaultClientParams = &virgil.VirgilClientParams{
		CardsServiceAddress:         conf.CardsServiceAddress,
		ReadonlyCardsServiceAddress: conf.ReadonlyCardsServiceAddress,
	}

	if len(conf.AppKey) == 0 {
		b, err := ioutil.ReadFile(conf.AppKeyPath)
		if err != nil {
			panic(err)
		}
		conf.AppKey = b
	}

	crypto := virgil.NewCrypto()
	appKey, err := crypto.ImportPrivateKey(conf.AppKey, conf.Passsword)
	if err != nil {
		panic(err)
	}

	return &Remote{
		client:     virgil.NewClient(token),
		appID:      conf.AppID,
		privateKey: appKey,
	}
}

type Remote struct {
	client     virgil.VirgilClient
	appID      string
	privateKey virgilcrypto.PrivateKey
}

func (s *Remote) GetCard(id string) (models.CardResponse, error) {
	var res models.CardResponse

	card, err := s.client.GetCard(id)
	if err != nil {
		return res, err
	}
	return mapCardToCardRequest(card), nil
}

func (s *Remote) SearchCards(c models.Criteria) (models.CardsResponse, error) {
	var res models.CardsResponse
	var scope enums.VirgilEnum

	if c.Scope == models.GlobalScope {
		scope = enums.CardScope.Global
	} else {
		scope = enums.CardScope.Application
	}

	cards, err := s.client.SearchCards(&search.Criteria{
		IdentityType: c.IdentityType,
		Identities:   c.Identities,
		Scope:        scope,
	})

	if err != nil {
		return res, err
	}

	for _, v := range cards {
		res = append(res, mapCardToCardRequest(v))
	}

	return res, nil
}

func (s *Remote) CreateCard(c models.CardResponse) (models.CardResponse, error) {

	// var scope enums.VirgilEnum
	// if c.Scope == "global" {
	// 	scope = enums.CardScope.Global
	// } else {
	// 	scope = enums.CardScope.Application
	// }

	// pk, err := virgilcrypto.DecodePublicKey(c.PublicKey)
	// if err != nil {
	// 	return models.CardResponse{}, err
	// }
	// ccr, err := virgil.NewCreateCardRequest(c.IdentityType, c.Identity, pk, scope, c.Data)
	// ccr.Signatures = c.Meta.Signatures

	// requestSigner := virgil.RequestSigner{
	// 	Crypto: virgil.NewCrypto(),
	// }
	// err = requestSigner.AuthoritySign(ccr, s.appID, s.privateKey)
	// if err != nil {
	// 	return models.CardResponse{}, err
	// }

	// card, err := s.client.CreateCard(ccr)
	// if err != nil {
	// 	return models.CardResponse{}, err
	// }
	// return mapCardToCardRequest(card), nil
	return c, nil
}

func mapCardToCardRequest(card *virgil.Card) models.CardResponse {
	return models.CardResponse{
		ID:       card.ID,
		Snapshot: card.Snapshot,
		Meta: models.ResponseMeta{
			CreatedAt:   card.CreatedAt,
			CardVersion: card.CardVersion,
			Signatures:  card.Signatures,
		},
	}
}
