package storage

import virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"

type CardStorage interface {
	GetCard(id string) (*virgil.Card, error)
	SearchCards(criteria virgil.Criteria) ([]*virgil.Card, error)
	CreateCard(request *virgil.SignableRequest) (*virgil.Card, error)
	RevokeCard(request *virgil.SignableRequest) error
}
