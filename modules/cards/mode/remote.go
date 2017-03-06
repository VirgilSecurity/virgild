package mode

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/pkg/errors"
	virgil "gopkg.in/virgil.v4"
)

type VirgilClient interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(*virgil.Criteria) ([]*virgil.Card, error)
	CreateCard(req *virgil.SignableRequest) (*virgil.Card, error)
	RevokeCard(req *virgil.SignableRequest) error
}

func vcard2Card(vcard *virgil.Card) *core.Card {
	c := &core.Card{
		ID:       vcard.ID,
		Snapshot: vcard.Snapshot,
		Meta: core.CardMeta{
			CreatedAt:   vcard.CreatedAt,
			CardVersion: vcard.CardVersion,
			Signatures:  vcard.Signatures,
			Relations:   vcard.Relations,
		},
	}
	if c.Meta.Relations == nil {
		c.Meta.Relations = make(map[string][]byte)
	}
	return c
}

type RemoteCardsMiddleware struct {
	Client VirgilClient
}

func (rcm *RemoteCardsMiddleware) Get(id string) (*core.Card, error) {
	vc, err := rcm.Client.GetCard(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return vcard2Card(vc), nil
}

func (rcm *RemoteCardsMiddleware) Search(criteria *virgil.Criteria) ([]core.Card, error) {
	vcards, err := rcm.Client.SearchCards(criteria)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cards := make([]core.Card, 0, len(vcards))
	for _, vc := range vcards {
		cards = append(cards, *vcard2Card(vc))
	}
	return cards, nil
}

func (rcm *RemoteCardsMiddleware) Create(req *core.CreateCardRequest) (*core.Card, error) {
	vcard, err := rcm.Client.CreateCard(&req.Request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return vcard2Card(vcard), nil
}

func (rcm *RemoteCardsMiddleware) Revoke(req *core.RevokeCardRequest) error {
	err := rcm.Client.RevokeCard(&req.Request)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
