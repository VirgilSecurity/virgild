package inegrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	url := fmt.Sprintf("%v/account/%v/collections/applications", c.Address, c.accountID)
	body, err := c.send(url)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetApplications [account_id: %s]", c.accountID)
	}

	var apps []core.Application
	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetApplications [unmarshal apps (account_id: %s)]", c.accountID)
	}
	return apps, nil
}

func (c *DevPortalClient) GetTokens() ([]core.Token, error) {
	url := fmt.Sprintf("%v/account/%v/collections/tokens", c.Address, c.accountID)
	body, err := c.send(url)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetTokens [account_id: %s]", c.accountID)
	}

	var tokens []core.Token
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.GetTokens [unmarshal tokens (account_id: %s)]", c.accountID)
	}
	return tokens, nil
}

func (c *DevPortalClient) Authorize(login, password string) error {
	doer := c.getDoer()
	body, _ := json.Marshal(map[string]string{
		"email":    login,
		"password": password,
	})
	req, _ := http.NewRequest(http.MethodPost, c.Address+"/authorization", bytes.NewReader(body))

	resp, err := doer.Do(req)
	if err != nil {
		return errors.Wrap(err, "DevPortalClient.Authorize [send request]")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var dpErr devPortalError
		err = json.NewDecoder(resp.Body).Decode(&dpErr)
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

func (c *DevPortalClient) send(url string) ([]byte, error) {
	if c.token == "" {
		return nil, errors.New("DevPortalClient.Send [not authorized]")
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", "bearer "+c.token)

	doer := c.getDoer()
	resp, err := doer.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DevPortalClient.Send [send request]")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "DevPortalClient.Send [read body]")
	}

	if resp.StatusCode != http.StatusOK {
		var dpErr devPortalError
		err = json.Unmarshal(body, &dpErr)
		if err != nil {
			return nil, errors.Wrapf(err, "DevPortalClient.Send [unmarshal err obj (status code: %v)]", resp.StatusCode)
		}
		return nil, dpErr
	}
	return body, nil
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
