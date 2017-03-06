package mode

import (
	"encoding/json"
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/pkg/errors"
	virgil "gopkg.in/virgil.v4"
)

type CardRepository interface {
	Get(id string) (*core.SqlCard, error)
	Find(identitis []string, identityType string, scope string) ([]core.SqlCard, error)
	Add(cs core.SqlCard) error
	MarkDeletedById(id string) error
	DeleteById(id string) error
	DeleteBySearch(identitis []string, identityType string, scope string) error
}

func sqlCard2Card(sql *core.SqlCard) (*core.Card, error) {
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

func card2SqlCard(card *core.Card) (*core.SqlCard, error) {
	var info virgil.CardModel
	err := json.Unmarshal(card.Snapshot, &info)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot read card snapshot")
	}

	jcard, err := json.Marshal(card)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot marshal card to json")
	}
	return &core.SqlCard{
		CardID:       card.ID,
		Identity:     info.Identity,
		IdentityType: info.IdentityType,
		Scope:        string(info.Scope),
		Card:         jcard,
	}, nil
}

type LocalCardsMiddleware struct {
	Repo CardRepository
}

func (lcm *LocalCardsMiddleware) Get(next core.GetCard) core.GetCard {
	getFromNext := func(id string) (card *core.Card, err error) {
		var scard *core.SqlCard

		card, err = next(id)
		if err != nil {
			return
		}
		scard, err = card2SqlCard(card)
		if err != nil {
			return
		}

		lcm.Repo.Add(*scard)
		return
	}

	return func(id string) (card *core.Card, err error) {
		var scard *core.SqlCard

		scard, err = lcm.Repo.Get(id)
		if err == core.ErrorEntityNotFound {
			return getFromNext(id)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "LocalCardsMiddleware get card(%v) local repository error", id)
		}

		if scard.ExpireAt < time.Now().Unix() {
			lcm.Repo.DeleteById(scard.CardID)
			return getFromNext(id)
		}

		if scard.Deleted {
			return nil, core.ErrorEntityNotFound
		}
		return sqlCard2Card(scard)
	}
}
func (lcm *LocalCardsMiddleware) Search(next core.SearchCards) core.SearchCards {
	serchFromNext := func(criteria *virgil.Criteria) ([]core.Card, error) {
		cards, err := next(criteria)
		if err != nil {
			return nil, err
		}

		for _, c := range cards {
			sqlCard, err := card2SqlCard(&c)
			if err != nil {
				return nil, err
			}
			lcm.Repo.Add(*sqlCard)
		}
		return cards, nil
	}

	return func(criteria *virgil.Criteria) ([]core.Card, error) {
		cards, err := lcm.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
		if err != nil {
			return nil, err
		}
		if len(cards) == 0 {
			return serchFromNext(criteria)
		}

		ct := time.Now().Unix()
		for _, v := range cards {
			if v.ExpireAt < ct {
				lcm.Repo.DeleteBySearch(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
				return serchFromNext(criteria)
			}
		}
		return sqlCards2Cards(cards)
	}
}

func (lcm *LocalCardsMiddleware) Create(next core.CreateCard) core.CreateCard {
	return func(req *core.CreateCardRequest) (card *core.Card, err error) {
		card, err = next(req)
		if err != nil {
			return
		}
		var scard *core.SqlCard
		scard, err = card2SqlCard(card)
		if err != nil {
			return
		}

		lcm.Repo.Add(*scard)

		return
	}
}

func (lcm *LocalCardsMiddleware) Revoke(next core.RevokeCard) core.RevokeCard {
	return func(req *core.RevokeCardRequest) (err error) {
		err = next(req)
		if err == nil {
			err = lcm.Repo.MarkDeletedById(req.Info.ID)
		}
		return
	}
}
