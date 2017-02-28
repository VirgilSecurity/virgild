package mode

import (
	"encoding/json"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/errors"
)

type VirgilClient interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(*virgil.Criteria) ([]*virgil.Card, error)
	CreateCard(req *virgil.SignableRequest) (*virgil.Card, error)
	RevokeCard(req *virgil.SignableRequest) error
}

type CardRepository interface {
	Get(id string) (*core.SqlCard, error)
	Find(identitis []string, identityType string, scope string) ([]core.SqlCard, error)
	Add(cs core.SqlCard) error
	DeleteById(id string) error
	DeleteBySearch(identitis []string, identityType string, scope string) error
}

func vcard2SqlCard(vcard *virgil.Card) (*core.SqlCard, error) {
	card := vcard2Card(vcard)
	jcard, err := json.Marshal(card)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &core.SqlCard{
		CardID:       vcard.ID,
		Identity:     vcard.Identity,
		IdentityType: vcard.IdentityType,
		Scope:        string(vcard.Scope),
		Card:         jcard,
	}, nil
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

func sqlCard2Card(sql *core.SqlCard) (*core.Card, error) {
	if sql.ErrorCode != 0 {
		return nil, core.ResponseErrorCode(sql.ErrorCode)
	}
	card := new(core.Card)
	err := json.Unmarshal(sql.Card, card)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return card, nil
}

func sqlCards2Cards(sql []core.SqlCard) ([]core.Card, error) {
	cards := make([]core.Card, 0)
	for _, v := range sql {
		c, err := sqlCard2Card(&v)
		if err == nil {
			cards = append(cards, *c)
		}
	}
	return cards, nil
}

func getFromRemote(remote VirgilClient, repo CardRepository, id string) (*core.Card, error) {
	vc, err := remote.GetCard(id)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			code := verr.ServiceErrorCode()
			if verr.HTTPErrorCode() == 404 {
				code = int(core.ErrorEntityNotFound)
				err = core.ErrorEntityNotFound
			}
			repo.Add(core.SqlCard{
				CardID:    id,
				ErrorCode: code,
			})
		}
		return nil, errors.Wrap(err, "")
	}
	sqlCard, err := vcard2SqlCard(vc)
	if err != nil {
		return nil, err
	}
	repo.Add(*sqlCard)
	return vcard2Card(vc), nil
}

func searchFromRemote(remote VirgilClient, repo CardRepository, criteria *virgil.Criteria) ([]core.Card, error) {
	vcards, err := remote.SearchCards(criteria)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			for k := range criteria.Identities {
				repo.Add(core.SqlCard{
					Identity:     criteria.Identities[k],
					IdentityType: criteria.IdentityType,
					Scope:        string(criteria.Scope),
					ErrorCode:    verr.ServiceErrorCode(),
				})
			}
		}
		return nil, errors.Wrap(err, "")
	}

	cards := make([]core.Card, 0)
	for _, vc := range vcards {

		sqlCard, err := vcard2SqlCard(vc)
		if err != nil {
			return nil, err
		}
		repo.Add(*sqlCard)

		cards = append(cards, *vcard2Card(vc))
	}
	return cards, nil
}
