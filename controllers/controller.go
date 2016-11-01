package controllers

import (
	"encoding/json"
	"errors"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (*models.CardResponse, error)
	SearchCards(models.Criteria) (models.CardsResponse, error)
	CreateCard(models.CardResponse) (*models.CardResponse, error)
	RevokeCard(id string, c models.CardResponse) error
}

type Controller struct {
	Storage Storage
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
		return nil, errors.New("Data has incorrect format")
	}

	cr.Scope = models.ResolveScope(cr.Scope)
	cards, err := c.Storage.SearchCards(cr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cards)
}

func (c *Controller) CreateCard(data []byte) ([]byte, error) {
	var cr models.CardResponse
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return nil, errors.New("Data has incorrect format")
	}
	card, err := c.Storage.CreateCard(cr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(card)
}

func (c *Controller) RevokeCard(id string, data []byte) error {
	var cr models.CardResponse
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return errors.New("Data has incorrect format")
	}
	return c.Storage.RevokeCard(id, cr)
}
