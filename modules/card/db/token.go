package db

import (
	"database/sql"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TokenStore struct {
	DB *sqlx.DB
}

func (s TokenStore) GetByValue(val string) (*core.Token, error) {
	token := new(core.Token)
	err := s.DB.Get(token, s.DB.Rebind("SELECT * from tokens where value=?"), val)
	if err == sql.ErrNoRows {
		return nil, core.EntityNotFoundErr
	}

	if err != nil {
		return nil, errors.Wrapf(err, "TokenStore.GetByVal(%s)", val)
	}
	return token, nil
}

func (s *TokenStore) Add(token core.Token) error {
	_, err := s.DB.NamedExec("INSERT into tokens(id, name, value, is_active, application_id, created_at, updated_at) VALUES(:id, :name, :value, :is_active, :application_id, :created_at, :updated_at)", token)
	if err != nil {
		return errors.Wrap(err, "TokenStore.Add")
	}
	return nil
}

func (s *TokenStore) Delete(id string) error {
	_, err := s.DB.Exec(s.DB.Rebind("DELETE tokens where id=?"), id)
	if err != nil {
		return errors.Wrap(err, "TokenStore.Delete")
	}
	return nil
}

func (s *TokenStore) GetAll() ([]core.Token, error) {
	var tokens []core.Token
	err := s.DB.Select(&tokens, "SELECT * from tokens")
	if err != nil {
		return tokens, errors.Wrap(err, "TokenStore.GetAll")
	}
	return tokens, nil
}
