package core

import (
	"time"

	virgil "gopkg.in/virgil.v4"
)

type Endpoint int

const (
	GetCardEndpoint     Endpoint = iota
	SearchCardsEndpoint Endpoint = iota
	CreateCardEndpoint  Endpoint = iota
	RevokeCardEndpoint  Endpoint = iota
)

type RequestStatistics struct {
	Data     time.Time
	Token    string
	Method   string
	Resource string
}

type SqlCard struct {
	CardID       string `xorm:"Index 'card_id'"`
	Identity     string `xorm:"Index"`
	IdentityType string
	Scope        string
	ExpireAt     time.Time
	Deleted      bool
	ErrorCode    int
	Card         []byte
}

type Card struct {
	ID       string   `json:"id"`
	Snapshot []byte   `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     CardMeta `json:"meta"`
}
type CardMeta struct {
	CreatedAt   string            `json:"created_at,omitempty"`
	CardVersion string            `json:"card_version,omitempty"`
	Signatures  map[string][]byte `json:"signs"`
}

type CreateCardRequest struct {
	Info    virgil.CardModel
	Request virgil.SignableRequest
}

type RevokeCardRequest struct {
	Info    virgil.RevokeCardRequest
	Request virgil.SignableRequest
}
