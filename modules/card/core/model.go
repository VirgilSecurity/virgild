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
