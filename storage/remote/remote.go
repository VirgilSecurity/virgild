package remote

import (
	"encoding/json"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/enums"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4/search"
)

type RemoteConfig struct {
	CardsServiceAddress         string
	ReadonlyCardsServiceAddress string
	PublicKey                   []byte
	AppID                       string
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

	crypto := virgil.NewCrypto()
	pk, err := crypto.ImportPublicKey(conf.PublicKey)
	if err != nil {
		panic(err)
	}
	v := virgil.NewCardsValidator(crypto)
	v.AddVerifier(conf.AppID, pk)

	client := virgil.NewClient(token)
	client.SetCardsValidator(v)

	return &Remote{
		client: virgil.NewClient(token),
	}
}

type Remote struct {
	client virgil.VirgilClient
}

func (s *Remote) GetCard(id string) (*models.CardResponse, error) {
	card, err := s.client.GetCard(id)
	if err != nil {
		return nil, err
	}
	return mapCardToCardRequest(card), nil
}

func (s *Remote) SearchCards(c models.Criteria) (models.CardsResponse, error) {
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
		return nil, err
	}

	res := models.CardsResponse{}
	for _, v := range cards {
		res = append(res, *mapCardToCardRequest(v))
	}

	return res, nil
}

func (s *Remote) CreateCard(c models.CardResponse) (*models.CardResponse, error) {
	vrs := virgil.SignedResponse{
		ID: c.ID,
		Meta: virgil.ResponseMeta{
			CreatedAt:   c.Meta.CreatedAt,
			CardVersion: c.Meta.CardVersion,
			Signatures:  c.Meta.Signatures,
		},
		Snapshot: c.Snapshot,
	}
	card, err := vrs.ToCard(virgil.NewCrypto())
	if err != nil {
		return nil, err
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

	card, err = s.client.CreateCard(r)
	if err != nil {
		return nil, err
	}
	return mapCardToCardRequest(card), nil
}

func (s *Remote) RevokeCard(id string, c models.CardResponse) error {
	r := virgil.NewEmptyRevokeCardRequest()

	json.Unmarshal(c.Snapshot, r)
	for k, v := range c.Meta.Signatures {
		r.AppendSignature(k, v)
	}

	return s.client.RevokeCard(r)
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
