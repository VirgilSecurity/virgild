package inegrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/pkg/errors"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type DevPortalClient struct {
	Doer      Doer
	Address   string
	token     string
	accountID string
}

func (c *DevPortalClient) GetApplications() ([]core.Application, error) {
	if c.token == "" {
		return nil, errors.New("DevPortalClient.GetApplications [not authorized]")
	}

	uri := fmt.Sprintf("%v/account/%v/collections/applications", c.Address, c.accountID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DevPortalClient.GetApplications [new request]")
	}
	req.Header.Add("Authorization", "bearer "+c.token)

	doer := c.getDoer()
	resp, err := doer.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetApplications [send request (account_id: %s)]", c.accountID)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apps []core.Application
		err = json.NewDecoder(resp.Body).Decode(&apps)
		if err != nil {
			return nil, errors.Wrapf(err, "DevPortalClient.GetApplications [unmarshal apps (account_id: %s)]", c.accountID)
		}
		return apps, nil
	}
	var dpErr devPortalError
	err = json.NewDecoder(resp.Body).Decode(&dpErr)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetApplications [unmarshal err obj (account_id: %s status code: %v)]", c.accountID, resp.StatusCode)
	}
	return nil, dpErr
}

func (c *DevPortalClient) GetTokens() ([]core.Token, error) {
	if c.token == "" {
		return nil, errors.New("DevPortalClient.GetTokens [not authorized]")
	}

	uri := fmt.Sprintf("%v/account/%v/collections/tokens", c.Address, c.accountID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DevPortalClient.GetTokens [new request]")
	}
	req.Header.Add("Authorization", "bearer "+c.token)

	doer := c.getDoer()
	resp, err := doer.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetTokens [send request (account_id: %s)]", c.accountID)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var tokens []core.Token
		err = json.NewDecoder(resp.Body).Decode(&tokens)
		if err != nil {
			return nil, errors.Wrapf(err, "DevPortalClient.GetTokens [unmarshal tokens (account_id: %s)]", c.accountID)
		}
		return tokens, nil
	}
	var dpErr devPortalError
	err = json.NewDecoder(resp.Body).Decode(&dpErr)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetTokens [unmarshal err obj (account_id: %s status code: %v)]", c.accountID, resp.StatusCode)
	}
	return nil, dpErr
}

func (c *DevPortalClient) Authorize(login, password string) error {
	doer := c.getDoer()
	body, _ := json.Marshal(map[string]string{
		"email":    login,
		"password": password,
	})
	req, err := http.NewRequest(http.MethodPost, c.Address+"/authorization", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "DevPortalClient.Authorize [new request]")
	}
	resp, err := doer.Do(req)
	if err != nil {
		return errors.Wrap(err, "DevPortalClient.Authorize [send request]")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var dpErr devPortalError
		err := json.NewDecoder(resp.Body).Decode(&dpErr)
		if err != nil {
			return errors.Wrapf(err, "DevPortalClient.Authorize [unmarshal error obj (resp status code: %v)]", resp.StatusCode)
		}
		return dpErr
	}
	var auth devPortalAuth
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return errors.Wrap(err, "DevPortalClient.Authorize [unmarshal auth info]")
	}
	c.accountID = auth.AccountID
	c.token = auth.AuthToken
	return nil
}

type devPortalAuth struct {
	AccountID string `json:"account_id"`
	AuthToken string `json:"auth_token"`
}

type devPortalError struct {
	Code    int
	Message string
}

func (err devPortalError) Error() string {
	return fmt.Sprintf("DevPortalError(code:%v msg: %v)", err.Code, err.Message)
}

func (c *DevPortalClient) getDoer() Doer {
	if c.Doer == nil {
		return http.DefaultClient
	}
	return c.Doer
}
