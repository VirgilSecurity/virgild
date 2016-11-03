package models

type Handler interface {
	GetCard(id string) (*models.CardResponse, *models.ErrorResponse)
	SearchCards(models.Criteria) ([]models.CardResponse, *models.ErrorResponse)
	CreateCard(*models.CardResponse) (*models.CardResponse, *models.ErrorResponse)
	RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse
}
