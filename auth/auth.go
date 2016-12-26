package auth

import (
	"encoding/json"
	"github.com/virgilsecurity/virgild/models"
)

type AuthHander struct {
	Token string
}

func (h *AuthHander) Auth(auth string) (bool, []byte) {
	if auth == ("VIRGIL " + h.Token) {
		return true, nil
	} else {
		jErr, _ := json.Marshal(models.MakeError(20300))
		return false, jErr
	}
}
