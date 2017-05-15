package background

import "github.com/VirgilSecurity/virgild/modules/card/core"

type FakeDevPortap struct {
}

func (dp *FakeDevPortap) GetApplications() ([]core.Application, error) {
	return []core.Application{
		core.Application{
			ID:          "1",
			Name:        "app1",
			Bundle:      "com.virgilsecurity.tochka.app1",
			CardID:      "1234",
			Description: "it's fake app",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
		},
		core.Application{
			ID:          "2",
			Name:        "app2",
			Bundle:      "com.virgilsecurity.tochka.app2",
			CardID:      "5678",
			Description: "it's fake app",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
		},
		core.Application{
			ID:          "3",
			Name:        "app3",
			Bundle:      "com.virgilsecurity.tochka.app3",
			CardID:      "4321",
			Description: "it's fake app",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
		},
	}, nil
}

func (dp *FakeDevPortap) GetTokens() ([]core.Token, error) {
	return []core.Token{
		core.Token{
			Name:        "my token",
			Value:       "1234",
			Active:      true,
			ID:          "1234",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
			Application: "1",
		},
		core.Token{
			Name:        "my token",
			Value:       "5678",
			Active:      true,
			ID:          "5678",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
			Application: "2",
		},
		core.Token{
			Name:        "my token",
			Value:       "4321",
			Active:      true,
			ID:          "4321",
			CreatedAt:   "2016-12-31T23:59:60+00:20",
			UpdatedAt:   "2016-12-31T23:59:60+00:20",
			Application: "3",
		},
	}, nil
}
