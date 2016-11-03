package protocols

type CodeResponse int

const (
	Ok           CodeResponse = iota
	RequestError CodeResponse = iota
	NotFound     CodeResponse = iota
	ServerError  CodeResponse = iota
)

type Controller interface {
	GetCard(id string) ([]byte, CodeResponse)
	SearchCards([]byte) ([]byte, CodeResponse)
	CreateCard([]byte) ([]byte, CodeResponse)
	RevokeCard(id string, data []byte) ([]byte, CodeResponse)
}
type AuthHandler interface {
	Auth(string) (bool, []byte)
}

type Server interface {
	Serve() error
}
