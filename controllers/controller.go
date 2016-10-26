package controllers

import (
	"encoding/json"
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (models.CardResponse, error)
	// Search(models.SearchCriteria)
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
