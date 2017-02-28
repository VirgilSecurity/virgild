package core

import virgil "gopkg.in/virgil.v4"

type SqlCard struct {
	CardID       string `xorm:"Index 'card_id'"`
	Identity     string `xorm:"Index"`
	IdentityType string
	Scope        string
	ExpireAt     int64
	Deleted      bool
	ErrorCode    int `xorm:"notnull"`
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
	Relations   map[string][]byte `json:"relations"`
}

type CreateCardRequest struct {
	Info    virgil.CardModel
	Request virgil.SignableRequest
}

type RevokeCardRequest struct {
	Info    virgil.RevokeCardRequest
	Request virgil.SignableRequest
}
