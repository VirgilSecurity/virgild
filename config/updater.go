package config

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Updater struct {
	app *App
}

func (u *Updater) Config() Config {
	c := u.app.Common.Config
	// clera private Data
	c.Cards.Signer.PrivateKey = ""
	c.Cards.Signer.PrivateKeyPassword = ""

	return c
}

func (u *Updater) Update(conf Config) error {
	err := u.validate(conf)
	if err != nil {
		return err
	}
	err = saveConfigToFole(conf, u.app.Common.ConfigPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("./restart.sh")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	go cmd.Run()

	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()

	return nil
}

func (u *Updater) validate(conf Config) error {
	ac := u.app.Common.Config
	if ac == conf {
		return fmt.Errorf("Configurations are equal")
	}
	if ac.DB != conf.DB {
		_, err := initDB(conf.DB)
		if err != nil {
			return err
		}
	}
	if ac.LogFile != conf.LogFile {
		_, err := initLogger(conf.LogFile)
		if err != nil {
			return err
		}
	}
	if ac.Cards.Remote != conf.Cards.Remote {
		_, err := initRemote(conf.Cards.Remote)
		if err != nil {
			return err
		}
	}
	if ac.Cards.Signer != conf.Cards.Signer {
		_, err := initSigner(&conf.Cards.Signer)
		if err != nil {
			return err
		}
	}
	if ac.Cards.VRA != conf.Cards.VRA {
		_, err := initVRA(conf.Cards.VRA)
		if err != nil {
			return err
		}
	}
	_, err := initAtuh(conf.Auth)
	if err != nil {
		return err
	}
	return nil
}
