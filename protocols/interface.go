package protocols

type CodeResponse int

const (
	Ok           = iota
	RequestError = iota
	NotFound     = iota
	ServerError  = iota
)

type Controller interface {
	GetCard(id string) ([]byte, CodeResponse)
	SearchCards([]byte) ([]byte, CodeResponse)
	CreateCard([]byte) ([]byte, CodeResponse)
	RevokeCard(id string, data []byte) CodeResponse
}
type AuthHandler interface {
	Auth(string) (bool, []byte)
}

type Server interface {
	Serve() error
}
