package mode

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

type AppModeCardHandler struct {
	Repo   CardRepository
	Remote VirgilClient
}

func (h *AppModeCardHandler) Get(id string) (*core.Card, error) {
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

func (h *AppModeCardHandler) Search(criteria *virgil.Criteria) ([]core.Card, error) {
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

func (h *AppModeCardHandler) Create(req *core.CreateCardRequest) (*core.Card, error) {
	vcard, err := h.Remote.CreateCard(&req.Request)
	if err != nil {
		return nil, err
	}
	sqlCard, err := vcard2SqlCard(vcard)
	if err != nil {
		return nil, err
	}

	h.Repo.Add(*sqlCard)
	return vcard2Card(vcard), nil
}

func (h *AppModeCardHandler) Revoke(req *core.RevokeCardRequest) error {
	err := h.Remote.RevokeCard(&req.Request)
	if err != nil {
		return err
	}

	h.Repo.DeleteById(req.Info.ID)
	return nil
}
