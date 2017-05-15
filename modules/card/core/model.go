package core

import virgil "gopkg.in/virgil.v4"

type Card struct {
	Response virgil.CardResponse
	Info     virgil.CardModel
}

type CreateCardRequest struct {
	Info    virgil.CardModel
	Request virgil.SignableRequest
}

type RevokeCardRequest struct {
	Info    virgil.RevokeCardRequest
	Request virgil.SignableRequest
}

type CreateRelationRequest struct {
	ID      string
	Request virgil.SignableRequest
}

type RevokeRelationRequest struct {
	ID      string
	Info    virgil.RevokeCardRequest
	Request virgil.SignableRequest
}

type Token struct {
	Name        string `json:"name" storm:"id" db:"name"`
	Value       string `json:"value" db:"value"`
	Active      bool   `json:"is_active" db:"is_active"`
	ID          string `json:"id" db:"id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	Application string `json:"application" db:"application_id"`
}

type Application struct {
	ID          string `json:"id" storm:"id" db:"id"`
	CardID      string `json:"card_id" db:"card_id"`
	Name        string `json:"name" db:"name"`
	Bundle      string `json:"bundle" db:"bundle"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}
