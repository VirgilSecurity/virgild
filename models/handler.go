package models

type Handler interface {
	GetCard(id string) (*CardResponse, *ErrorResponse)
	SearchCards(Criteria) ([]CardResponse, *ErrorResponse)
	CreateCard(*CardResponse) (*CardResponse, *ErrorResponse)
	RevokeCard(id string, c *CardResponse) *ErrorResponse
}
