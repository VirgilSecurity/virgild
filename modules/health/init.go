package health

import (
	"github.com/VirgilSecurity/virgild/config"
)

func Init(app *config.App) *HealthChecker {
	return &HealthChecker{}
}
