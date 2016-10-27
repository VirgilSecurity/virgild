package controllers

import (
	"encoding/json"
	"errors"
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (models.CardResponse, error)
	SearchCards(models.Criteria) (models.CardsResponse, error)
}

type Controller struct {
	Storage Storage
}

func (c *Controller) GetCard(id string) ([]byte, error) {
	card, err := c.Storage.GetCard(id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(card)
}

func (c *Controller) SearchCards(data []byte) ([]byte, error) {
	var cr models.Criteria
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return nil, errors.New("Data has incorrect format")
	}

	if cr.Scope == "" {
		cr.Scope = models.ApplicationScope
	}

	cards, err := c.Storage.SearchCards(cr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cards)
}
