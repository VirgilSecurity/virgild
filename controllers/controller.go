package controllers

import (
	"encoding/json"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (*models.CardResponse, error)
	SearchCards(models.Criteria) ([]models.CardResponse, error)
	CreateCard(*models.CardResponse) (*models.CardResponse, error)
	RevokeCard(id string, c *models.CardResponse) error
}

type Validator interface {
	Validate(*models.CardResponse) error
}

type Controller struct {
	Storage   Storage
	Validator Validator
}

func (c *Controller) GetCard(id string) ([]byte, error) {
	card, err := c.Storage.GetCard(id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, nil
	}
	return json.Marshal(card)
}

func (c *Controller) SearchCards(data []byte) ([]byte, error) {
	var cr models.Criteria
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return nil, models.ErrorResponse{
			Code: 30000,
		}
	}

	cr.Scope = models.ResolveScope(cr.Scope)
	cards, err := c.Storage.SearchCards(cr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cards)
}

func (c *Controller) CreateCard(data []byte) ([]byte, error) {
	cr := new(models.CardResponse)
	err := json.Unmarshal(data, cr)
	if err != nil {
		return nil, models.ErrorResponse{
			Code: 30000,
		}
	}
	err = c.Validator.Validate(cr)
	if err != nil {
		return nil, err
	}
	card, err := c.Storage.CreateCard(cr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(card)
}

func (c *Controller) RevokeCard(id string, data []byte) error {
	cr := new(models.CardResponse)
	err := json.Unmarshal(data, cr)
	if err != nil {
		return models.ErrorResponse{
			Code: 30000,
		}
	}

	err = c.Validator.Validate(cr)
	if err != nil {
		return err
	}
	return c.Storage.RevokeCard(id, cr)
}
