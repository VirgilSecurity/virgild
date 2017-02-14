package mode

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

type CacheModeHandler struct {
	Repo   CardRepository
	Remote VirgilClient
}

func (h *CacheModeHandler) Get(id string) (*core.Card, error) {
	c, err := h.Repo.Get(id)
	if err == core.ErrorEntityNotFound {
		return getFromRemote(h.Remote, h.Repo, id)
	}
	if err != nil {
		return nil, err
	}
	if c.ExpireAt < time.Now().Unix() {
		h.Repo.DeleteById(id)
		return getFromRemote(h.Remote, h.Repo, id)
	}

	return sqlCard2Card(c)
}

func (h *CacheModeHandler) Search(criteria *virgil.Criteria) ([]core.Card, error) {
	cards, err := h.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return searchFromRemote(h.Remote, h.Repo, criteria)
	}

	for _, v := range cards {
		if v.ExpireAt < time.Now().Unix() {
			h.Repo.DeleteBySearch(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
			return searchFromRemote(h.Remote, h.Repo, criteria)
		}
	}
	return sqlCards2Cards(cards)
}

func (h *CacheModeHandler) Create(req *core.CreateCardRequest) (*core.Card, error) {
	return nil, core.ErrorForbidden
}

func (h *CacheModeHandler) Revoke(req *core.RevokeCardRequest) error {
	return core.ErrorForbidden
}
