package remote

import (
	"encoding/json"
	"fmt"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
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

	client := virgil.NewClient(token)

	return &Remote{
		client: client,
	}
}

type Remote struct {
	client virgil.VirgilClient
}

func (s *Remote) GetCard(id string) (*models.CardResponse, *models.ErrorResponse) {
	card, err := s.client.GetCard(id)
	if err != nil {
		return nil, models.MakeError(10000)
	}
	return mapCardToCardRequest(card), nil
}

func (s *Remote) SearchCards(c models.Criteria) ([]models.CardResponse, *models.ErrorResponse) {
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
		return nil, models.MakeError(10000)
	}

	res := []models.CardResponse{}
	for _, v := range cards {
		res = append(res, *mapCardToCardRequest(v))
	}

	return res, nil
}

func (s *Remote) CreateCard(c *models.CardResponse) (*models.CardResponse, *models.ErrorResponse) {
	vrs := virgil.SignedResponse{
		ID: c.ID,
		Meta: virgil.ResponseMeta{
			CreatedAt:   c.Meta.CreatedAt,
			CardVersion: c.Meta.CardVersion,
			Signatures:  c.Meta.Signatures,
		},
		Snapshot: c.Snapshot,
	}
	card, err := vrs.ToCard()
	jCard, _ := json.MarshalIndent(card, "", "\t")
	fmt.Println("Restored card:", string(jCard[:]))
	if err != nil {
		return nil, models.MakeError(10000)
	}
	r := virgil.NewEmptyCreateCardRequest()
	r.Data = card.Data
	r.DeviceInfo = card.DeviceInfo
	r.Identity = card.Identity
	r.IdentityType = card.IdentityType
	r.PublicKey, _ = card.PublicKey.Encode()
	r.Scope = card.Scope
	for k, v := range card.Signatures {
		r.AppendSignature(k, v)
	}
	jR, _ := json.MarshalIndent(card, "", "\t")
	fmt.Println("New Card request:", string(jR[:]))

	card, err = s.client.CreateCard(r)
	if err != nil {
		return nil, models.MakeError(10000)
	}
	return mapCardToCardRequest(card), nil
}

func (s *Remote) RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse {
	r := virgil.NewEmptyRevokeCardRequest()

	json.Unmarshal(c.Snapshot, r)
	for k, v := range c.Meta.Signatures {
		r.AppendSignature(k, v)
	}
	err := s.client.RevokeCard(r)
	if err != nil {
		return models.MakeError(10000)
	} else {
		return nil
	}
}

func mapCardToCardRequest(card *virgil.Card) *models.CardResponse {
	return &models.CardResponse{
		ID:       card.ID,
		Snapshot: card.Snapshot,
		Meta: models.ResponseMeta{
			CreatedAt:   card.CreatedAt,
			CardVersion: card.CardVersion,
			Signatures:  card.Signatures,
		},
	}
}
