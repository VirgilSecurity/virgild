package models

type CardResponse struct {
	ID       string       `json:"id,omitempty"`
	Snapshot []byte       `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	CreatedAt   string            `json:"created_at,omitempty"`
	CardVersion string            `json:"card_version,omitempty"`
	Signatures  map[string][]byte `json:"signs"`
}
