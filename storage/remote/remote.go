package remote

import (
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/enums"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/search"
)

type RemoteConfig struct {
	CardsServiceAddress         string
	ReadonlyCardsServiceAddress string
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

	return &Remote{
		client: virgil.NewClient(token),
	}
}

type Remote struct {
	client virgil.VirgilClient
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

func (s *Remote) CreateCard(c models.CardRequst) (models.CardResponse, error) {
	appID := "[YOUR_APP_ID_HERE]"
	appKeyPassword := "[YOUR_APP_KEY_PASSWORD_HERE]"
	appKeyData, err := ioutil.ReadFile("[YOUR_APP_KEY_PATH_HERE]")
	appKey, err := crypto.ImportPrivateKey(appKeyData, appKeyPassword)

	ccr := virgil.CreateCardRequest{
		Identity:     c.Identity,
		IdentityType: c.IdentityType,
		PublicKey:    c.PublicKey,
		Scope:        c.Scope,
		SignableRequest: virgil.SignableRequest{
			Signatures: models.Signatures,
		},
		Signatures: c.Signatures,
	}
	err = requestSigner.AuthoritySign(ccr, appID, appKey)
	return virgil.CreateCard(ccr)
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