package mode

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/errors"
)

type VirgilClient interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(virgil.Criteria) ([]*virgil.Card, error)
	CreateCard(req *virgil.SignableRequest) (*virgil.Card, error)
	RevokeCard(req *virgil.SignableRequest) error
}

type AppModeCardHandler struct {
	Repo   CardRepository
	Remote VirgilClient
}

func (h *AppModeCardHandler) remoteGet(id string) (*core.Card, error) {
	vc, err := h.Remote.GetCard(id)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			code := verr.ServiceErrorCode()
			if verr.HTTPErrorCode() == 404 {
				code = int(core.ErrorEntityNotFound)
				err = core.ErrorEntityNotFound
			}
			h.Repo.Add(core.SqlCard{
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
	h.Repo.Add(*sqlCard)
	return vcard2Card(vc), nil
}

func (h *AppModeCardHandler) Get(id string) (*core.Card, error) {
	c, err := h.Repo.Get(id)
	if err == core.ErrorEntityNotFound {
		return h.remoteGet(id)
	}
	if err != nil {
		return nil, err
	}
	if c.ExpireAt < time.Now().Unix() {
		h.Repo.DeleteById(id)
		return h.remoteGet(id)
	}

	return sqlCard2Card(c)
}

func (h *AppModeCardHandler) remoteSearch(criteria *virgil.Criteria) ([]core.Card, error) {
	vcards, err := h.Remote.SearchCards(*criteria)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			for k := range criteria.Identities {
				h.Repo.Add(core.SqlCard{
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
		h.Repo.Add(*sqlCard)

		cards = append(cards, *vcard2Card(vc))
	}
	return cards, nil
}

func (h *AppModeCardHandler) Search(criteria *virgil.Criteria) ([]core.Card, error) {
	cards, err := h.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return h.remoteSearch(criteria)
	}

	for _, v := range cards {
		if v.ExpireAt < time.Now().Unix() {
			h.Repo.DeleteBySearch(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
			return h.remoteSearch(criteria)
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
