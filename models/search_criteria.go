package models

const (
	GlobalScope      string = "global"
	ApplicationScope string = "application"
)

type Criteria struct {
	Scope        string   `json:"scope,omitempty"`
	IdentityType string   `json:"identity_type,omitempty"`
	Identities   []string `json:"identities"`
}

func ResolveScope(scope string) string {
	if scope == GlobalScope {
		return GlobalScope
	} else {
		return ApplicationScope
	}
}
