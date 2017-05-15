package background

import (
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/pkg/errors"
)

type AppStore interface {
	Add(app core.Application) error
	Delete(id string) error
	GetAll() ([]core.Application, error)
}

type TokenStore interface {
	Add(token core.Token) error
	Delete(id string) error
	GetAll() ([]core.Token, error)
}

type DevPortalClient interface {
	GetApplications() ([]core.Application, error)
	GetTokens() ([]core.Token, error)
}

func UpdateAppsCronJob(appStore AppStore, client DevPortalClient) func() error {
	return func() error {
		devApps, err := client.GetApplications()
		if err != nil {
			return errors.Wrap(err, "UpdateApplications.GetDevPortalApplications")
		}

		locApps, err := appStore.GetAll()
		if err != nil {
			return errors.Wrap(err, "UpdateApplications.GetAll")
		}

		var addApps []core.Application
		for _, devApp := range devApps {
			var exist = false
			for i, locApp := range locApps {
				if devApp.ID == locApp.ID && devApp.UpdatedAt == locApp.UpdatedAt {
					locApps = append(locApps[:i], locApps[i+1:]...)
					exist = true
					break
				}
			}
			if !exist {
				addApps = append(addApps, devApp)
			}
		}

		for _, delApp := range locApps {
			err := appStore.Delete(delApp.ID)
			if err != nil {
				return errors.Wrapf(err, "UpdateApplications.Delete(%s)", delApp.ID)
			}
		}

		for _, addApp := range addApps {
			err := appStore.Add(addApp)
			if err != nil {
				return errors.Wrapf(err, "UpdateApplications.Add(%s)", addApp.ID)
			}
		}

		return nil
	}
}

func UpdateTokensCronJob(tokenStore TokenStore, client DevPortalClient) func() error {
	return func() error {
		devTokens, err := client.GetTokens()
		if err != nil {
			return errors.Wrap(err, "UpdateTokens.GetDevPortalTokens")
		}
		locTokens, err := tokenStore.GetAll()
		if err != nil {
			return errors.Wrap(err, "UpdateTokens.GetAll")
		}
		var addTokens []core.Token

		for _, devToken := range devTokens {
			var exist = false
			for i, locToken := range locTokens {
				if devToken.ID == locToken.ID && devToken.UpdatedAt == locToken.UpdatedAt {
					locTokens = append(locTokens[:i], locTokens[i+1:]...)
					exist = true
					break
				}
			}
			if !exist {
				addTokens = append(addTokens, devToken)
			}
		}

		for _, delToken := range locTokens {
			err := tokenStore.Delete(delToken.ID)
			if err != nil {
				return errors.Wrapf(err, "UpdateTokens.Delete(%s)", delToken.ID)
			}
		}
		for _, addToken := range addTokens {
			err := tokenStore.Add(addToken)
			if err != nil {
				return errors.Wrapf(err, "UpdateTokens.Add(%s)", addToken.ID)
			}
		}

		return nil
	}
}
