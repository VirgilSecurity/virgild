package core

type Request struct {
	Snapshot []byte      `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     RequestMeta `json:"meta"`
}

type RequestMeta struct {
	Signatures map[string][]byte `json:"signs"`
}

type Criteria struct {
	Scope        string   `json:"scope,omitempty"`
	IdentityType string   `json:"indentity_type,omitempty"`
	Identities   []string `json:"identities"`
}
