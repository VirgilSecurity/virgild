package sync

import (
	"github.com/virgilsecurity/virgild/models"
)

type Storage interface {
	GetCard(id string) (*models.CardResponse, *models.ErrorResponse)
	SearchCards(models.Criteria) ([]models.CardResponse, *models.ErrorResponse)
	CreateCard(*models.CardResponse) (*models.CardResponse, *models.ErrorResponse)
	RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse
}

type Sync struct {
	Local  Storage
	Remote Storage
}

func (s Sync) GetCard(id string) (*models.CardResponse, *models.ErrorResponse) {
	c, err := s.Local.GetCard(id)
	if err != nil {
		return s.Remote.GetCard(id)
	}
	if c == nil {
		c, err = s.Remote.GetCard(id)
		if c != nil {
			s.Local.CreateCard(c)
		}
	}
	return c, err
}

func (s Sync) SearchCards(c models.Criteria) ([]models.CardResponse, *models.ErrorResponse) {
	csl, err := s.Local.SearchCards(c)
	if err != nil {
		return s.Remote.SearchCards(c)
	}

	// TODO identity can be assign to multiple cards
	if len(csl) != len(c.Identities) {
		csr, err := s.Remote.SearchCards(c)
		if err != nil {
			return csr, err
		}

		for _, vr := range csr {
			var exist bool = false
			for _, vl := range csl {
				if vr.ID == vl.ID {
					exist = true
					break
				}
			}
			if !exist {
				s.Local.CreateCard(&vr)
			}
		}

		return csr, nil
	}
	return csl, nil
}

func (s Sync) CreateCard(c *models.CardResponse) (*models.CardResponse, *models.ErrorResponse) {
	r, err := s.Remote.CreateCard(c)
	if err != nil {
		return nil, err
	}
	s.Local.CreateCard(r)
	return r, nil
}

func (s Sync) RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse {
	err := s.Remote.RevokeCard(id, c)
	if err != nil {
		return err
	}
	s.Local.RevokeCard(id, c)
	return nil
}
