package health

import (
	"time"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/go-xorm/xorm"
)

func Init(app *config.App) *HealthChecker {

	return &HealthChecker{
		checkList: map[string]info{
			"db":            makeDBHealth(app.Common.DB),
			"cards-service": makeCardsServiceHealth(app.Cards.Remote.VClient, app.Common.Config.Cards.Remote.Authority.CardID),
		},
	}
}

func makeDBHealth(db *xorm.Engine) info {
	return func() (map[string]interface{}, error) {
		start := time.Now()
		err := db.Ping()
		end := time.Now()

		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"latency": end.Sub(start) / time.Millisecond,
		}, nil
	}
}

func makeCardsServiceHealth(c *virgil.Client, id string) info {
	return func() (map[string]interface{}, error) {
		start := time.Now()
		_, err := c.GetCard(id)
		end := time.Now()

		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"latency": end.Sub(start) / time.Millisecond,
		}, nil
	}
}
