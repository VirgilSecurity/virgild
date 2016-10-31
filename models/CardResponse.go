package models

type CardsResponse []CardResponse

type CardResponse struct {
	ID       string       `json:"id"`
	Snapshot []byte       `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	CreatedAt   string            `json:"created_at"`
	CardVersion string            `json:"card_version"`
	Signatures  map[string][]byte `json:"signs"`
}
