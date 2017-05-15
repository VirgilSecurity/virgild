package inegrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/VirgilSecurity/virgild/modules/card/integrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeDoer struct {
	mock.Mock
}

func (f *fakeDoer) Do(req *http.Request) (resp *http.Response, err error) {
	args := f.Called(req)
	resp, _ = args.Get(0).(*http.Response)
	err = args.Error(1)

	return
}

func makeAuthDevPortalClient(token string, accountID string) *inegrations.DevPortalClient {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"account_id":"%s","auth_token":"%s"}`, accountID, token))),
	}

	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c := &inegrations.DevPortalClient{Address: "http://address", Doer: doer}

	c.Authorize("login", "password")
	return c
}

func TestDevPortalClient_ClientIsNotAuthorize_ReturnErr(t *testing.T) {
	c := &inegrations.DevPortalClient{Address: "http://localhost"}

	_, err := c.GetApplications()
	assert.Error(t, err)

	_, err = c.GetTokens()
	assert.Error(t, err)
}

func TestDevPortalClientAuthorize_DoerReturnErr_ReturnErr(t *testing.T) {
	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(nil, fmt.Errorf("Error"))
	c := &inegrations.DevPortalClient{Address: "http://localhost", Doer: doer}
	err := c.Authorize("login", "password")
	assert.Error(t, err)
}

func TestDevPortalClientAuthorize_AuthorizeDevErr_ReturnErr(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"code":123,"message":"code error"}`)),
	}

	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c := &inegrations.DevPortalClient{Address: "http://localhost", Doer: doer}
	err := c.Authorize("login", "password")
	assert.Error(t, err)
}

func TestDevPortalClientAuthorize_UnsupportedBodyResp_ReturnErr(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"bad"body data"}`)),
	}

	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c := &inegrations.DevPortalClient{Address: "http://localhost", Doer: doer}
	err := c.Authorize("login", "password")
	assert.Error(t, err)
}

func TestDevPortalClientAuthorize_AuthorizeSeccessUnsupportedBody_ReturnNil(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"bad"body data"}`)),
	}

	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c := &inegrations.DevPortalClient{Address: "http://address", Doer: doer}

	err := c.Authorize("login", "password")

	assert.Error(t, err)
}

func TestDevPortalClientAuthorize_AuthorizeSeccess_ReturnNil(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"account_id":"1234","auth_token":"12344"}`)),
	}

	matchFunc := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodPost &&
			req.URL.String() == "http://address/authorization" &&
			bytes.Compare(body, []byte(`{"email":"login","password":"password"}`)) == 0
	}
	doer := new(fakeDoer)
	doer.On("Do", mock.MatchedBy(matchFunc)).Return(resp, nil)
	c := &inegrations.DevPortalClient{Address: "http://address", Doer: doer}

	err := c.Authorize("login", "password")

	assert.NoError(t, err)
}

func TestDevPortalClientTable_DoerReturnErr_ReturnErr(t *testing.T) {
	table := []func(c *inegrations.DevPortalClient) error{
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetApplications()
			return err
		},
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetTokens()
			return err
		},
	}

	c := makeAuthDevPortalClient("token", "acc_id")
	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(nil, fmt.Errorf("Error"))
	c.Doer = doer

	for _, f := range table {
		err := f(c)
		assert.Error(t, err)
	}
}

func TestDevPortalClientTable_DevPortalReturnErr_ReturnErr(t *testing.T) {
	table := []func(c *inegrations.DevPortalClient) error{
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetApplications()
			return err
		},
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetTokens()
			return err
		},
	}

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"code":123,"message":"code error"}`)),
	}
	c := makeAuthDevPortalClient("token", "acc_id")
	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c.Doer = doer

	for _, f := range table {
		err := f(c)
		assert.Error(t, err)
	}
}

func TestDevPortalClientTable_UnsupportedBodyErrResp_ReturnErr(t *testing.T) {
	table := []func(c *inegrations.DevPortalClient) error{
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetApplications()
			return err
		},
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetTokens()
			return err
		},
	}

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"bad"body data"}`)),
	}
	c := makeAuthDevPortalClient("token", "acc_id")
	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c.Doer = doer

	for _, f := range table {
		err := f(c)
		assert.Error(t, err)
	}
}

func TestDevPortalClientTable_SeccessUnsupportedBody_ReturnNil(t *testing.T) {
	table := []func(c *inegrations.DevPortalClient) error{
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetApplications()
			return err
		},
		func(c *inegrations.DevPortalClient) error {
			_, err := c.GetTokens()
			return err
		},
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"bad"body data"}`)),
	}
	c := makeAuthDevPortalClient("token", "acc_id")
	doer := new(fakeDoer)
	doer.On("Do", mock.Anything).Return(resp, nil)
	c.Doer = doer

	for _, f := range table {
		err := f(c)
		assert.Error(t, err)
	}
}

func TestDevPortalClientGetTokens_Seccess_ReturnVal(t *testing.T) {
	expected := []core.Token{
		core.Token{
			ID:          "token id 1",
			Name:        "test token",
			Active:      false,
			Value:       "token value",
			CreatedAt:   "created date",
			UpdatedAt:   "updated date",
			Application: "applications",
		},
		core.Token{
			ID:          "token id 2",
			Name:        "test token",
			Active:      false,
			Value:       "token value",
			CreatedAt:   "created date",
			UpdatedAt:   "updated date",
			Application: "applications",
		},
	}
	body, _ := json.Marshal(expected)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	c := makeAuthDevPortalClient("token", "acc_id")
	matchFunc := func(req *http.Request) bool {
		return req.Method == http.MethodGet &&
			req.URL.String() == "http://address/account/acc_id/collections/tokens"
	}
	doer := new(fakeDoer)
	doer.On("Do", mock.MatchedBy(matchFunc)).Return(resp, nil)
	c.Doer = doer
	c.Address = "http://address"

	tokens, err := c.GetTokens()

	assert.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestDevPortalClientGetApplications_Seccess_ReturnVal(t *testing.T) {
	expected := []core.Application{
		core.Application{
			ID:          "token id 1",
			Name:        "test token",
			CardID:      "application card id",
			Bundle:      "com.virgilsecurity.application",
			Description: "application description",
			CreatedAt:   "created date",
			UpdatedAt:   "updated date",
		},
		core.Application{
			ID:          "token id 2",
			Name:        "test token",
			CardID:      "application card id",
			Bundle:      "com.virgilsecurity.application",
			Description: "application description",
			CreatedAt:   "created date",
			UpdatedAt:   "updated date",
		},
	}
	body, _ := json.Marshal(expected)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	c := makeAuthDevPortalClient("token", "acc_id")
	matchFunc := func(req *http.Request) bool {
		return req.Method == http.MethodGet &&
			req.URL.String() == "http://address/account/acc_id/collections/applications"
	}
	doer := new(fakeDoer)
	doer.On("Do", mock.MatchedBy(matchFunc)).Return(resp, nil)
	c.Doer = doer
	c.Address = "http://address"

	apps, err := c.GetApplications()

	assert.NoError(t, err)
	assert.Equal(t, expected, apps)
}
