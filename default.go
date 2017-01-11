package main

import (
	"time"

	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type DefaultModeCardHandler struct {
	Repo        CardRepository
	Signer      RequestSigner
	Validator   Validator
	Fingerprint Fingerprint
}

func (h *DefaultModeCardHandler) Get(id string) (interface{}, error) {
	card, err := h.Repo.Get(id)
	if err != nil {
		return nil, err
	}
	return sqlCard2Card(card)
}
func (h *DefaultModeCardHandler) Search(criteria *virgil.Criteria) (interface{}, error) {
	if ok, err := h.Validator.IsValidSearchCriteria(criteria); !ok {
		return nil, err
	}
	cards, err := h.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
	if err != nil {
		return nil, err
	}
	return sqlCards2Cards(cards)
}
func (h *DefaultModeCardHandler) Create(req *CreateCardRequest) (interface{}, error) {
	if ok, err := h.Validator.IsValidCreateCardRequest(req); !ok {
		return nil, err
	}

	err := h.Signer.Sign(&req.Request)
	if err != nil {
		return nil, err
	}

	vcard := &virgil.Card{
		ID:           h.Fingerprint.Calculate(req.Request.Snapshot),
		Snapshot:     req.Request.Snapshot,
		Identity:     req.Info.Identity,
		IdentityType: req.Info.IdentityType,
		Scope:        req.Info.Scope,
		CreatedAt:    time.Now().Format(time.RFC3339),
		CardVersion:  "v4",
		Signatures:   req.Request.Signatures,
		Data:         req.Info.Data,
		DeviceInfo:   req.Info.DeviceInfo,
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

func (h *DefaultModeCardHandler) Revoke(req *RevokeCardRequest) (interface{}, error) {
	if ok, err := h.Validator.IsValidRevokeCardRequest(req); !ok {
		return nil, err
	}
	return nil, h.Repo.DeleteById(req.Info.ID)
}
