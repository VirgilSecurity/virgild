package sync

import (
	"fmt"
	"github.com/VirgilSecurity/virgil-apps-cards-cacher/models"
)

type Storage interface {
	GetCard(id string) (models.CardResponse, error)
	SearchCards(models.Criteria) (models.CardsResponse, error)
	CreateCard(models.CardResponse) (models.CardResponse, error)
	RevokCard(id string, c models.CardResponse) error
}

type Sync struct {
	Local  Storage
	Remote Storage
}

func (s Sync) GetCard(id string) (models.CardResponse, error) {
	c, err := s.Local.GetCard(id)
	if err != nil {
		fmt.Println("Miss cache")
		c, err = s.Remote.GetCard(id)
		if err != nil {
			return c, err
		}
		return s.Local.CreateCard(c)
	}
	return c, err
}

func (s Sync) SearchCards(c models.Criteria) (models.CardsResponse, error) {

	csl, err := s.Local.SearchCards(c)
	if err != nil {
		return csl, err
	}

	if len(csl) != len(c.Identities) {
		fmt.Println("Miss cache local:", len(csl), "try find:", len(c.Identities))
		csr, err := s.Remote.SearchCards(c)
		if err != nil {
			return csr, err
		}

		for _, vr := range csr {
			var exist bool = false
			fmt.Println("remote:", vr.ID)
			for _, vl := range csl {
				fmt.Println("local:", vl.ID)
				if vr.ID == vl.ID {
					exist = true
					break
				}
			}
			fmt.Println("Exist:", exist)
			if !exist {
				s.Local.CreateCard(vr)
			}
		}
		return csr, nil
	}
	return csl, nil
}

func (s Sync) CreateCard(c models.CardResponse) (models.CardResponse, error) {
	r, err := s.Remote.CreateCard(c)
	if err != nil {
		return models.CardResponse{}, err
	}
	_, err = s.Local.CreateCard(r)
	if err != nil {
		fmt.Printf("Local storage err:", err)
	}
	return r, nil
}

func (s Sync) RevokCard(id string, c models.CardResponse) error {
	err := s.Remote.RevokCard(id, c)
	if err != nil {
		return err
	}
	err = s.Local.RevokCard(id, c)
	if err != nil {
		return err
	}
	return nil
}
