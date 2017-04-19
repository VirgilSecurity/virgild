package core

import (
	"context"

	virgil "gopkg.in/virgil.v4"
)

type GetCardHandler func(ctx context.Context, id string) (*virgil.CardResponse, error)
type SearchCardsHandler func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error)
type CreateCardHandler func(ctx context.Context, createReq *CreateCardRequest) (*virgil.CardResponse, error)
type RevokeCardHandler func(ctx context.Context, revokeReq *RevokeCardRequest) error
type CreateRelationHandler func(ctx context.Context, createReq *CreateRelationRequest) (*virgil.CardResponse, error)
type RevokeRelationHandler func(ctx context.Context, revokeReq *RevokeRelationRequest) (*virgil.CardResponse, error)
