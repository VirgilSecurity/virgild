package storage

import (
	"errors"
	"fmt"
	virgil "gopkg.in/virgilsecurity/virgil-sdk-go.v4"
	"testing"
)

func Test_ServiceGet_SearchInStorageCard(t *testing.T) {
	expected := virgil.Card{ID: "test"}

	storage := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return &virgil.Card{ID: id}, nil
		},
	}
	service := MakeService(&storage, nil, nil)
	actual, _ := service.Get(expected.ID)

	if actual.ID != expected.ID {
		t.Fatal("Service doesn't search in storage:", expected.ID, "expected but", actual.ID, "actual")
	}
}

func Test_ServiceGet_SearchInRemouteStorageCard(t *testing.T) {
	expected := virgil.Card{ID: "test"}

	storage := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return nil, nil
		},
	}

	remoute := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return &virgil.Card{ID: id}, nil
		},
	}

	service := MakeService(&storage, &remoute, nil)
	actual, _ := service.Get(expected.ID)

	if actual.ID != expected.ID {
		t.Fatal("Service doesn't search in remoute storage:", expected.ID, "expected but", actual.ID, "actual")
	}
}

func Test_ServiceGet_LoggedIfSomethingWrongInStorage(t *testing.T) {
	storage := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return nil, errors.New("error")
		},
	}
	remoute := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return &virgil.Card{ID: id}, nil
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Get("test")

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceGet_LoggedIfSomethingWrongInRemouteStorage(t *testing.T) {
	storage := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return nil, nil
		},
	}
	remoute := fakeCardService{
		GetFunc: func(id string) (*virgil.Card, error) {
			return nil, errors.New("errror")
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Get("test")

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceFind_SearchInStorageCard(t *testing.T) {
	storage := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return []*virgil.Card{&virgil.Card{}}, nil
		},
	}
	service := MakeService(&storage, nil, nil)
	service.Find("email", "test")

	if !storage.FindInvoked {
		t.Fatal("Service doesn't search in storage")
	}
}

func Test_ServiceFind_SearchInRemouteStorageCard(t *testing.T) {
	storage := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return nil, nil
		},
	}

	remoute := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return []*virgil.Card{&virgil.Card{}}, nil
		},
	}

	service := MakeService(&storage, &remoute, nil)
	service.Find("email", "test")

	if !remoute.FindInvoked {
		t.Fatal("Service doesn't search in remoute storage")
	}
}

func Test_ServiceFind_LoggedIfSomethingWrongInStorage(t *testing.T) {
	storage := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return nil, errors.New("error")
		},
	}
	remoute := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return []*virgil.Card{&virgil.Card{}}, nil
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Find("email", "test")

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceFind_LoggedIfSomethingWrongInRemouteStorage(t *testing.T) {
	storage := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return nil, nil
		},
	}
	remoute := fakeCardService{
		FindFunc: func(identityType string, identities ...string) ([]*virgil.Card, error) {
			return nil, errors.New("errror")
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Find("email", "test")

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceCreate_SearchInStorageCard(t *testing.T) {
	storage := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}
	remoute := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}

	service := MakeService(&storage, &remoute, nil)
	service.Create(virgil.Card{})

	if !storage.CreateInvoked {
		t.Fatal("Service doesn't search in storage")
	}
}

func Test_ServiceCreate_SearchInRemouteStorageCard(t *testing.T) {
	storage := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}

	remoute := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}

	service := MakeService(&storage, &remoute, nil)
	service.Create(virgil.Card{})

	if !remoute.CreateInvoked {
		t.Fatal("Service doesn't search in remoute storage")
	}
}

func Test_ServiceCreate_LoggedIfSomethingWrongInStorage(t *testing.T) {
	storage := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return errors.New("error")
		},
	}
	remoute := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Create(virgil.Card{})

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceCreate_LoggedIfSomethingWrongInRemouteStorage(t *testing.T) {
	storage := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return nil
		},
	}
	remoute := fakeCardService{
		CreateFunc: func(Card virgil.Card) error {
			return errors.New("errror")
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Create(virgil.Card{})

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceRevoke_SearchInStorageCard(t *testing.T) {
	storage := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return nil
		},
	}
	remoute := fakeCardService{
		RevokeFunc: func(Card virgil.Card) error {
			return nil
		},
	}
	service := MakeService(&storage, &remoute, nil)
	service.Revoke(virgil.Card{})

	if !storage.RevokeInvoked {
		t.Fatal("Service doesn't search in storage")
	}
}

func Test_ServiceRevoke_SearchInRemouteStorageCard(t *testing.T) {
	storage := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return nil
		},
	}

	remoute := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return nil
		},
	}

	service := MakeService(&storage, &remoute, nil)
	service.Revoke(virgil.Card{})

	if !remoute.RevokeInvoked {
		t.Fatal("Service doesn't search in remoute storage")
	}
}

func Test_ServiceRevoke_LoggedIfSomethingWrongInStorage(t *testing.T) {
	storage := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return errors.New("error")
		},
	}
	remoute := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return nil
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Revoke(virgil.Card{})

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

func Test_ServiceRevoke_LoggedIfSomethingWrongInRemouteStorage(t *testing.T) {
	storage := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return nil
		},
	}
	remoute := fakeCardService{
		RevokeFunc: func(card virgil.Card) error {
			return errors.New("errror")
		},
	}
	var log fakeLogger

	service := MakeService(&storage, &remoute, &log)
	service.Revoke(virgil.Card{})

	if !log.IsInvoked {
		t.Error("Service doesn't loging a error")
	}
}

type fakeLogger struct {
	IsInvoked bool
	ErrorText string
}

func (f *fakeLogger) Printf(format string, v ...interface{}) {
	f.IsInvoked = true
	f.ErrorText = fmt.Sprintf(format, v...)
}

func (f *fakeLogger) Println(v ...interface{}) {
	f.IsInvoked = true
	f.ErrorText = fmt.Sprintln(v...)
}

type fakeCardService struct {
	GetInvoked    bool
	GetFunc       func(string) (*virgil.Card, error)
	FindInvoked   bool
	FindFunc      func(identityType string, identities ...string) ([]*virgil.Card, error)
	CreateInvoked bool
	CreateFunc    func(virgil.Card) error
	RevokeInvoked bool
	RevokeFunc    func(virgil.Card) error
}

func (f *fakeCardService) Get(id string) (*virgil.Card, error) {
	f.GetInvoked = true
	return f.GetFunc(id)
}

func (f *fakeCardService) Find(identityType string, identities ...string) ([]*virgil.Card, error) {
	f.FindInvoked = true
	return f.FindFunc(identityType, identities...)
}

func (f *fakeCardService) Create(card virgil.Card) error {
	f.CreateInvoked = true
	return f.CreateFunc(card)
}

func (f *fakeCardService) Revoke(card virgil.Card) error {
	f.RevokeInvoked = true
	return f.RevokeFunc(card)
}
