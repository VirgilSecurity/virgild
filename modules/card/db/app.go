package db

import (
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ApplicationsStore struct {
	DB *sqlx.DB
}

// func (s *ApplicationsStore) GetById(id string) (*core.Application, error) {
// 	app := new(core.Application)
// 	err := s.DB.Get(app, s.DB.Rebind("SELECT * from applications where id=?"), id)
// 	if err == sql.ErrNoRows {
// 		return nil, core.EntityNotFoundErr
// 	}
//
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "ApplicationsStore.GetById(%s)", id)
// 	}
// 	return app, nil
// }

func (s *ApplicationsStore) Add(app core.Application) error {
	_, err := s.DB.NamedExec("INSERT into applications(id, card_id,name ,bundle,description,created_at ,updated_at) VALUES(:id,:card_id,:name,:bundle,:description, :created_at, :updated_at)", app)
	if err != nil {
		return errors.Wrap(err, "ApplicationsStore.Add")
	}
	return nil
}

func (s *ApplicationsStore) Delete(id string) error {
	_, err := s.DB.Exec(s.DB.Rebind("DELETE applications where id=?"), id)
	if err != nil {
		return errors.Wrap(err, "ApplicationsStore.Delete")
	}
	return nil
}

func (s *ApplicationsStore) GetAll() ([]core.Application, error) {
	var apps []core.Application
	err := s.DB.Select(&apps, "SELECT * from applications")
	if err != nil {
		return apps, errors.Wrap(err, "ApplicationsStore.GetAll")
	}
	return apps, nil
}
