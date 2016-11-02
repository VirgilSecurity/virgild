package local

import (
	"encoding/json"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
)

type CardRepository interface {
	Get(id string) (*CardSql, error)
	Find(models.Criteria) ([]CardSql, error)
	Add(CardSql) error
	Delete(string) error
}

type CardRequest struct {
	Identity     string            `json:"identity"`
	IdentityType string            `json:"identity_type"`
	PublicKey    []byte            `json:"public_key"` //DER encoded public key
	Scope        string            `json:"scope"`
	Data         map[string]string `json:"data,omitempty"`
	DeviceInfo   DeviceInfo        `json:"info"`
}

type DeviceInfo struct {
	Device     string `json:"device"`
	DeviceName string `json:"device_name"`
}

type Local struct {
	Repo CardRepository
}

func (s Local) GetCard(id string) (*models.CardResponse, error) {
	c, err := s.Repo.Get(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	r := new(models.CardResponse)
	err = json.Unmarshal([]byte(c.Card), &r)
	return r, err
}

func (s Local) SearchCards(c models.Criteria) ([]models.CardResponse, error) {
	var r []models.CardResponse
	cs, err := s.Repo.Find(c)
	if err != nil {
		return r, err
	}
	for _, v := range cs {
		var cr models.CardResponse
		err = json.Unmarshal([]byte(v.Card), &cr)
		if err != nil {
			return r, err
		}
		r = append(r, cr)
	}
	return r, err
}

func (s Local) CreateCard(c *models.CardResponse) (*models.CardResponse, error) {
	var cr CardRequest
	err := json.Unmarshal(c.Snapshot, &cr)
	if err != nil {
		return nil, models.ErrorResponse{Code: 30107}
	}

	jCard, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	cs := CardSql{
		Id:           c.ID,
		Identity:     cr.Identity,
		IdentityType: cr.IdentityType,
		Scope:        cr.Scope,
		Card:         string(jCard[:]),
	}
	err = s.Repo.Add(cs)
	return c, err
}

func (s Local) RevokeCard(id string, c *models.CardResponse) error {
	err := s.Repo.Delete(id)
	return err
}
