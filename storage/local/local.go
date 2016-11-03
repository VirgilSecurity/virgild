package local

import (
	"encoding/hex"
	"encoding/json"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type CardRepository interface {
	Get(id string) (*CardSql, error)
	Find(models.Criteria) ([]CardSql, error)
	Add(CardSql) error
	Delete(string) error
}

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
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
	Repo   CardRepository
	Logger Logger
}

func (s *Local) GetCard(id string) (*models.CardResponse, *models.ErrorResponse) {
	c, err := s.Repo.Get(id)
	if err != nil {
		s.Logger.Printf("Local storage [GetCard(%v)]: %s", id, err)
		return nil, models.MakeError(10000)
	}
	if c == nil {
		return nil, nil
	}
	r := new(models.CardResponse)
	err = json.Unmarshal([]byte(c.Card), r)
	if err != nil {
		s.Logger.Printf("Local storage [GetCard(%v)]: %s", id, err)
		return nil, models.MakeError(10000)
	}
	return r, nil
}

func (s *Local) SearchCards(c models.Criteria) ([]models.CardResponse, *models.ErrorResponse) {
	var r []models.CardResponse
	cs, err := s.Repo.Find(c)
	if err != nil {
		jc, _ := json.MarshalIndent(c, "", "\t")
		s.Logger.Printf("Local storage [SearchCard(%s)]: %s", jc, err)
		return nil, models.MakeError(10000)
	}
	for _, v := range cs {
		var cr models.CardResponse
		err = json.Unmarshal([]byte(v.Card), &cr)
		if err != nil {
			jc, _ := json.MarshalIndent(c, "", "\t")
			s.Logger.Printf("Local storage [SearchCard(%s)] on the value %s: %s", jc, v.Card, err)
			return nil, models.MakeError(10000)
		}
		r = append(r, cr)
	}
	return r, nil
}

func (s *Local) CreateCard(c *models.CardResponse) (*models.CardResponse, *models.ErrorResponse) {
	var cr CardRequest
	err := json.Unmarshal(c.Snapshot, &cr)
	if err != nil {
		jc, _ := json.MarshalIndent(c, "", "\t")
		s.Logger.Printf("Local storage [CreateCard(%s)]: %s", jc, err)
		return nil, models.MakeError(30107)
	}
	id := c.ID
	if id == "" {
		crypto := virgil.Crypto()
		fp := crypto.CalculateFingerprint(c.Snapshot)
		id = hex.EncodeToString(fp)
		c.ID = id
	}

	jCard, err := json.Marshal(c)
	if err != nil {
		jc, _ := json.MarshalIndent(c, "", "\t")
		s.Logger.Printf("Local storage [CreateCard(%s)]: %s", jc, err)
		return nil, models.MakeError(10000)
	}
	cs := CardSql{
		Id:           id,
		Identity:     cr.Identity,
		IdentityType: cr.IdentityType,
		Scope:        cr.Scope,
		Card:         string(jCard[:]),
	}
	err = s.Repo.Add(cs)
	if err != nil {
		jc, _ := json.MarshalIndent(c, "", "\t")
		s.Logger.Printf("Local storage [CreateCard(%s)]: %s", jc, err)
		return nil, models.MakeError(10000)
	}
	return c, nil
}

func (s *Local) RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse {
	err := s.Repo.Delete(id)
	if err != nil {
		jc, _ := json.MarshalIndent(c, "", "\t")
		s.Logger.Printf("Local storage [CreateCard(%v,%s)]: %s", id, jc, err)
		return models.MakeError(10000)
	}
	return nil
}
