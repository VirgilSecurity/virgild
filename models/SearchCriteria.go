package models

const (
	GlobalScope string = "global"
	Application string = "application"
)

type Criteria struct {
	Scope        string   `json:"scope"`
	IdentityType string   `json:"identity_type"`
	Identities   []string `json:"identities"`
}
