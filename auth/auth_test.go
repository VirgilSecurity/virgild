package auth

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"testing"
)

func Test_Validate_PrefixInvalid_ReturnError(t *testing.T) {
	expected := models.MakeError(20300)
	token := "123"
	auth := AuthHander{
		Token: token,
	}
	valided, err := auth.Auth("VV " + token)
	assert.False(t, valided)
	actual := new(models.ErrorResponse)
	json.Unmarshal(err, actual)
	assert.Equal(t, expected, actual)
}

func Test_Validate_TokenInvalid_ReturnError(t *testing.T) {
	expected := models.MakeError(20300)
	auth := AuthHander{
		Token: "123",
	}
	valided, err := auth.Auth("VIRGIL not valided")
	assert.False(t, valided)
	actual := new(models.ErrorResponse)
	json.Unmarshal(err, actual)
	assert.Equal(t, expected, actual)
}

func Test_Validate_TokenValid_ReturnNil(t *testing.T) {
	token := "123"
	auth := AuthHander{
		Token: token,
	}
	valided, err := auth.Auth("VIRGIL " + token)
	assert.True(t, valided)
	assert.Nil(t, err)
}
