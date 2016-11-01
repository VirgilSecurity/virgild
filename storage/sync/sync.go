package sync

import (
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (*models.CardResponse, error)
	SearchCards(models.Criteria) ([]models.CardResponse, error)
	CreateCard(*models.CardResponse) (*models.CardResponse, error)
	RevokeCard(id string, c *models.CardResponse) error
}

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type Sync struct {
	Logger Logger
	Local  Storage
	Remote Storage
}

func (s Sync) GetCard(id string) (*models.CardResponse, error) {
	c, err := s.Local.GetCard(id)
	if err != nil {
		s.Logger.Println("Local storage (GetCard):", err)
	}
	if c == nil {
		c, err = s.Remote.GetCard(id)
		if err != nil || c == nil {
			return nil, err
		}
		return s.Local.CreateCard(c)
	}
	return c, nil
}

func (s Sync) SearchCards(c models.Criteria) ([]models.CardResponse, error) {

	csl, err := s.Local.SearchCards(c)
	if err != nil {
		s.Logger.Println("Local storage (SearchCards):", err)
	}

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

func (s Sync) CreateCard(c *models.CardResponse) (*models.CardResponse, error) {
	r, err := s.Remote.CreateCard(c)
	if err != nil {
		return nil, err
	}
	_, err = s.Local.CreateCard(r)
	if err != nil {
		s.Logger.Println("Local storage (CreateCard):", err)
	}
	return r, nil
}

func (s Sync) RevokeCard(id string, c *models.CardResponse) error {
	err := s.Remote.RevokeCard(id, c)
	if err != nil {
		return err
	}
	err = s.Local.RevokeCard(id, c)
	if err != nil {
		s.Logger.Println("Local storage (RevokeCard):", err)
	}
	return nil
}
