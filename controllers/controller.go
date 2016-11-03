package controllers

import (
	"encoding/json"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
)

type Storage interface {
	GetCard(id string) (*models.CardResponse, *models.ErrorResponse)
	SearchCards(models.Criteria) ([]models.CardResponse, *models.ErrorResponse)
	CreateCard(*models.CardResponse) (*models.CardResponse, *models.ErrorResponse)
	RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse
}

type ServiceSigner interface {
	Sign(*models.CardResponse) error
}

type Controller struct {
	Storage Storage
	Signer  ServiceSigner
}

func (c *Controller) GetCard(id string) ([]byte, protocols.CodeResponse) {
	card, err := c.Storage.GetCard(id)
	if err != nil {
		return mapErrResponseToCodeResponse(err)
	}
	if card == nil {
		return nil, protocols.NotFound
	}
	jcard, _ := json.Marshal(card)
	return jcard, protocols.Ok
}

func (c *Controller) SearchCards(data []byte) ([]byte, protocols.CodeResponse) {
	var cr models.Criteria
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return mapErrResponseToCodeResponse(models.MakeError(30000))
	}

	cr.Scope = models.ResolveScope(cr.Scope)
	cards, e := c.Storage.SearchCards(cr)
	if e != nil {
		return mapErrResponseToCodeResponse(e)
	}
	jCards, _ := json.Marshal(cards)
	return jCards, protocols.Ok
}

func (c *Controller) CreateCard(data []byte) ([]byte, protocols.CodeResponse) {
	cr := new(models.CardResponse)
	err := json.Unmarshal(data, cr)
	if err != nil {
		return mapErrResponseToCodeResponse(models.MakeError(30000))
	}

	err = c.Signer.Sign(cr)
	if err != nil {
		return mapErrResponseToCodeResponse(models.MakeError(10000))
	}
	card, e := c.Storage.CreateCard(cr)
	if e != nil {
		return mapErrResponseToCodeResponse(e)
	}
	jCard, _ := json.Marshal(card)
	return jCard, protocols.Ok
}

func (c *Controller) RevokeCard(id string, data []byte) ([]byte, protocols.CodeResponse) {
	cr := new(models.CardResponse)
	err := json.Unmarshal(data, cr)
	if err != nil {
		return mapErrResponseToCodeResponse(models.MakeError(30000))
	}
	e := c.Storage.RevokeCard(id, cr)
	if e != nil {
		return mapErrResponseToCodeResponse(e)
	}
	return nil, protocols.Ok
}

// You can find list of actial code response by the following link
// https://virgilsecurity.com/docs/services/cards/v4.0(latest)/cards-service#appendix-a-response-codes
func mapErrResponseToCodeResponse(err *models.ErrorResponse) ([]byte, protocols.CodeResponse) {
	switch err.Code {
	case 10000:
		r, _ := json.Marshal(err)
		return r, protocols.ServerError
	default:
		r, _ := json.Marshal(err)
		return r, protocols.RequestError
	}
}
