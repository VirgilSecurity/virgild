package card

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/pkg/errors"
	virgil "gopkg.in/virgil.v4"
)

type cacheCardMiddleware struct {
	cache coreapi.Cache
}

func (c *cacheCardMiddleware) GetCard(f core.GetCardHandler) core.GetCardHandler {
	return func(ctx context.Context, id string) (card *virgil.CardResponse, err error) {
		owner := core.GetOwnerRequest(ctx)
		key := getCardKey(owner, id)
		has := c.cache.Get(key, &card)

		if has {
			return card, err
		}

		card, err = f(ctx, id)
		if err == nil {
			c.cache.Set(key, card)
		}

		return card, err
	}
}

func (c *cacheCardMiddleware) SearchCards(f core.SearchCardsHandler) core.SearchCardsHandler {
	return func(ctx context.Context, crit *virgil.Criteria) (cards []virgil.CardResponse, err error) {
		owner := core.GetOwnerRequest(ctx)
		if crit.Scope == virgil.CardScope.Global {
			owner = ""
		}

		p := []string{owner, crit.IdentityType, string(crit.Scope)}
		sort.Strings(crit.Identities)
		p = append(p, crit.Identities...)
		key := strings.Join(p, "_")

		var ids []string
		has := c.cache.Get(key, &ids)

		if has {
			cachePass := true
			for _, id := range ids {
				var card *virgil.CardResponse
				has = c.cache.Get(getCardKey(owner, id), &card)
				if !has {
					cachePass = false
					break
				}
				cards = append(cards, *card)
			}
			if cachePass {
				return cards, nil
			}
		}

		cards, err = f(ctx, crit)
		if err != nil {
			return nil, errors.Wrap(err, "CacheSerchCards(send)")
		}

		for _, card := range cards {
			c.cache.Set(getCardKey(owner, card.ID), card)
			ids = append(ids, card.ID)
		}

		c.cache.Set(key, ids)

		return cards, nil
	}
}

func (c cacheCardMiddleware) CreateCard(f core.CreateCardHandler) core.CreateCardHandler {
	return func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		card, err := f(ctx, req)
		if err != nil {
			return nil, errors.Wrap(err, "Cache.CreateCard(send)")
		}
		key := getCardKey(core.GetOwnerRequest(ctx), card.ID)
		c.cache.Set(key, card)

		return card, err
	}
}

func (c cacheCardMiddleware) RevokeCard(f core.RevokeCardHandler) core.RevokeCardHandler {
	return func(ctx context.Context, req *core.RevokeCardRequest) error {
		err := f(ctx, req)
		if err != nil {
			return errors.Wrap(err, "Cache.RevokeCard(send)")
		}

		key := getCardKey(core.GetOwnerRequest(ctx), req.Info.ID)
		c.cache.Del(key)
		return nil
	}
}

func (c cacheCardMiddleware) CreateRelations(f core.CreateRelationHandler) core.CreateRelationHandler {
	return func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {

		card, err := f(ctx, req)
		if err != nil {
			return nil, errors.Wrap(err, "Cache.CreateRelations(send)")
		}
		key := getCardKey(core.GetOwnerRequest(ctx), card.ID)
		c.cache.Set(key, card)
		return card, nil
	}
}

func (c cacheCardMiddleware) RevokeRelations(f core.RevokeRelationHandler) core.RevokeRelationHandler {
	return func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {

		card, err := f(ctx, req)
		if err != nil {
			return nil, errors.Wrap(err, "Cache.RevokeRelations(send)")
		}
		key := getCardKey(core.GetOwnerRequest(ctx), card.ID)
		c.cache.Set(key, card)
		return card, nil
	}
}

func getCardKey(owner string, id string) string {
	return fmt.Sprintf("%v_%v", owner, id)
}
