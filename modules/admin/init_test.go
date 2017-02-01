package admin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type fakeConfigRepo struct {
	mock.Mock
}

func (f *fakeConfigRepo) Config() config.Config {
	args := f.Called()
	return args.Get(0).(config.Config)
}
func (f *fakeConfigRepo) Update(conf config.Config) error {
	args := f.Called(conf)
	return args.Error(0)
}

func TestGetVirgilDCardInfo_ReturnVal(t *testing.T) {
	expected := config.VirgilDCard{
		CardID:    "123",
		PublicKey: "1234",
	}

	ctx := fasthttp.RequestCtx{}
	s := getVirgilDCardInfo(expected)
	s(&ctx)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())

	var actual config.VirgilDCard
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestUpdateConf_JSONInvalid_ReturnErr(t *testing.T) {
	type respErr struct {
		Message string `json:"message"`
	}
	s := updateConf(nil)

	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetBody([]byte(``))
	s(&ctx)

	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode())

	var actual respErr
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, "JSON invalid", actual.Message)
}

func TestUpdateConf_ConfigUpdaterReturnErr_ReturnErr(t *testing.T) {
	type respErr struct {
		Message string `json:"message"`
	}
	const errMsg = "Error"
	repo := new(fakeConfigRepo)
	repo.On("Update", mock.Anything).Return(fmt.Errorf(errMsg))

	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetBody([]byte(`{"common":{"db":"1234"}}`))
	s := updateConf(repo)
	s(&ctx)

	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode())

	var actual respErr
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, errMsg, actual.Message)
}

func TestUpdateConf_Seccess(t *testing.T) {
	expected := config.Config{
		Cards: config.CardsConfig{
			Signer: config.SignerConfig{
				CardID: "1234",
			},
		},
	}
	repo := new(fakeConfigRepo)
	repo.On("Update", expected).Return(nil)

	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetBody([]byte(`{"cards":{"signer":{"card_id":"1234"}}}`))
	s := updateConf(repo)
	s(&ctx)

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())

}

func TestGetConf_Return(t *testing.T) {
	expected := config.Config{
		Cards: config.CardsConfig{
			Signer: config.SignerConfig{
				CardID: "1234",
			},
		},
	}
	repo := new(fakeConfigRepo)
	repo.On("Config").Return(expected)

	ctx := fasthttp.RequestCtx{}
	s := getConf(repo)
	s(&ctx)

	assert.Equal(t, "application/json", string(ctx.Response.Header.ContentType()))
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())

	var actual config.Config
	err := json.Unmarshal(ctx.Response.Body(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
