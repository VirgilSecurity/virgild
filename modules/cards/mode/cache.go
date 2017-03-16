package mode

import (
	"fmt"
	"sort"
	"strings"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

type CacheCardsMiddleware struct {
	Manager config.CacheManager
}

func (ccm *CacheCardsMiddleware) Get(next core.GetCard) core.GetCard {
	return func(id string) (card *core.Card, err error) {
		if has := ccm.Manager.Get(id, &card); !has {
			card, err = next(id)
			if err != nil {
				return
			}
			ccm.Manager.Set(id, card)
			return
		}
		return
	}
}

func (ccm *CacheCardsMiddleware) Search(next core.SearchCards) core.SearchCards {
	return func(crit *virgil.Criteria) (cards []core.Card, err error) {
		var ids []string

		sort.Strings(crit.Identities)
		key := fmt.Sprint(crit.IdentityType, crit.Scope, strings.Join(crit.Identities, "_"))

		if ccm.Manager.Get(key, &ids) {
			var miss bool
			for _, id := range ids {
				var card core.Card
				if has := ccm.Manager.Get(id, &card); !has {
					miss = true
					break
				}
				cards = append(cards, card)
			}

			if !miss {
				return
			}
		}

		cards, err = next(crit)
		if err != nil {
			return
		}

		ids = ids[:0]
		for _, card := range cards {
			ids = append(ids, card.ID)
			ccm.Manager.Set(card.ID, &card)
		}

		ccm.Manager.Set(key, ids)
		return
	}
}

func (ccm *CacheCardsMiddleware) Create(next core.CreateCard) core.CreateCard {
	return func(req *core.CreateCardRequest) (card *core.Card, err error) {
		card, err = next(req)
		if err == nil {
			ccm.Manager.Set(card.ID, card)
		}
		return
	}
}

func (ccm *CacheCardsMiddleware) Revoke(next core.RevokeCard) core.RevokeCard {
	return func(req *core.RevokeCardRequest) (err error) {
		err = next(req)
		if err == nil {
			ccm.Manager.Del(req.Info.ID)
		}
		return
	}
}
