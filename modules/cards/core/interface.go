package core

import "gopkg.in/virgil.v4"

type GetCard func(id string) (*Card, error)
type SearchCards func(c *virgil.Criteria) ([]Card, error)
type CreateCard func(req *CreateCardRequest) (*Card, error)
type RevokeCard func(req *RevokeCardRequest) error
