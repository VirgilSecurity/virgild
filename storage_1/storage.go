package storage

import (
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	search "gopkg.in/virgilsecurity/virgil-sdk-go.v4/search"
)

type CardStorage interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(criteria *search.Criteria) ([]*virgil.Card, error)
}

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Service struct {
	storage CardService
	remoute CardService
	log     Logger
}

type CardService interface {
	Get(id string) (*virgil.Card, error)
	Find(identityType string, identities ...string) ([]*virgil.Card, error)
	Create(virgil.Card) error
	Revoke(virgil.Card) error
}

func MakeService(storage CardService, remoute CardService, log Logger) Service {
	return Service{
		storage: storage,
		remoute: remoute,
		log:     log,
	}
}

func (s *Service) Get(id string) (*virgil.Card, error) {
	card, err := s.storage.Get(id)
	if err != nil {
		s.log.Printf("[GET] Storage service return error: %s", err)
	}

	if card == nil {
		card, err = s.remoute.Get(id)
		if err != nil {
			s.log.Printf("[GET] Remoute storage service return error: %s", err)
		}
	}
	return card, err
}

func (s *Service) Find(identityType string, identities ...string) ([]*virgil.Card, error) {
	cards, err := s.storage.Find(identityType, identities...)
	if err != nil {
		s.log.Printf("[FIND] Storage service return error: %s", err)
	}

	if len(cards) == 0 {
		cards, err = s.remoute.Find(identityType, identities...)
		if err != nil {
			s.log.Printf("[FIND] Remoute storage service return error: %s", err)
		}
	}
	return cards, err
}

func (s *Service) Create(card virgil.Card) error {
	err := s.storage.Create(card)
	if err != nil {
		s.log.Printf("[CREATE] Storage service return error: %s", err)
	}

	err = s.remoute.Create(card)
	if err != nil {
		s.log.Printf("[CREATE] Remoute storage service return error: %s", err)
	}
	return err
}

func (s *Service) Revoke(card virgil.Card) error {
	err := s.storage.Revoke(card)
	if err != nil {
		s.log.Printf("[REVOKE] Storage service return error: %s", err)
	}

	err = s.remoute.Revoke(card)
	if err != nil {
		s.log.Printf("[REVOKE] Remoute storage service return error: %s", err)
	}
	return err
}
