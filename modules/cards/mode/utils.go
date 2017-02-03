package mode

import (
	"encoding/json"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/pkg/errors"
	virgil "gopkg.in/virgil.v4"
)

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
	return &core.Card{
		ID:       vcard.ID,
		Snapshot: vcard.Snapshot,
		Meta: core.CardMeta{
			CreatedAt:   vcard.CreatedAt,
			CardVersion: vcard.CardVersion,
			Signatures:  vcard.Signatures,
		},
	}
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
