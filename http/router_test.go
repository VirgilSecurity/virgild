package http

import (
	"github.com/stretchr/testify/mock"
)

type MockController struct {
	mock.Mock
}

func (c MockController) GetCard(id string) ([]byte, error) {
	args := c.Called(id)
	return args.Get(0).([]byte), args.Error(1)
}

func (c MockController) SearchCards(data []byte) ([]byte, error) {
	args := c.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (c MockController) CreateCard(data []byte) ([]byte, error) {
	args := c.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (c MockController) RevokeCard(id string, data []byte) error {
	args := c.Called(id, data)
	return args.Error(1)
}
