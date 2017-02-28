package mode

import (
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

type Fingerprint interface {
	Calculate(data []byte) string
}

type DefaultModeCardHandler struct {
	Repo        CardRepository
	Fingerprint Fingerprint
}

func (h *DefaultModeCardHandler) Get(id string) (*core.Card, error) {
	card, err := h.Repo.Get(id)
	if err != nil {
		return nil, err
	}
	return sqlCard2Card(card)
}
func (h *DefaultModeCardHandler) Search(criteria *virgil.Criteria) ([]core.Card, error) {
	cards, err := h.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
	if err != nil {
		return nil, err
	}
	return sqlCards2Cards(cards)
}
func (h *DefaultModeCardHandler) Create(req *core.CreateCardRequest) (*core.Card, error) {
	vcard := &virgil.Card{
		ID:           h.Fingerprint.Calculate(req.Request.Snapshot),
		Snapshot:     req.Request.Snapshot,
		Identity:     req.Info.Identity,
		IdentityType: req.Info.IdentityType,
		Scope:        req.Info.Scope,
		CreatedAt:    time.Now().Format(time.RFC3339),
		CardVersion:  "v4",
		Signatures:   req.Request.Meta.Signatures,
		Data:         req.Info.Data,
		DeviceInfo:   req.Info.DeviceInfo,
		Relations:    make(map[string][]byte, 0),
	}

	sqlCard, err := vcard2SqlCard(vcard)
	if err != nil {
		return nil, err
	}
	err = h.Repo.Add(*sqlCard)
	if err != nil {
		return nil, err
	}
	return vcard2Card(vcard), nil
}

func (h *DefaultModeCardHandler) Revoke(req *core.RevokeCardRequest) error {
	return h.Repo.DeleteById(req.Info.ID)
}
