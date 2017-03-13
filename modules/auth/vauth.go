package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type vclient struct {
	HttpClient HttpClient
	Host       string
}

type tokenInfo struct {
	OwnerID string
	Scope   []string
}

type vautherr struct {
	Code int `json:"code"`
}

func (c *vclient) Verify(token string) (*tokenInfo, error) {
	b := fmt.Sprintf(`{"access_token": "%v"}`, token)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/authorization/actions/verify", c.Host), strings.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	switch resp.StatusCode {
	case http.StatusBadRequest:
		e := new(vautherr)
		err = d.Decode(e)
		if err != nil {
			return nil, err
		}
		if e.Code == 53080 { //https://github.com/VirgilSecurity/virgil-services-auth
			return nil, errAuthServiceDenny
		}
		return nil, fmt.Errorf("Unexpected response for Virgil Auth (status: %v err_code: %v token: %v)", resp.StatusCode, e.Code, token)
	case http.StatusOK:
		info := new(tokenInfo)
		err = d.Decode(info)
		if err != nil {
			return nil, err
		}
		return info, nil
	default:
		return nil, fmt.Errorf("Unexpected response for Virgil Auth (status: %v token: %v)", resp.StatusCode, token)
	}
}

type vauthclient interface {
	Verify(t string) (*tokenInfo, error)
}

func externalScopes(c vauthclient) func(token string) ([]string, error) {
	return func(token string) ([]string, error) {
		info, err := c.Verify(token)
		if err != nil {
			return make([]string, 0), err
		}
		return info.Scope, errTokenInvalid
	}
}
