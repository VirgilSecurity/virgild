package core

type SymmetricKey struct {
	KeyID        string `xorm:"key_id" json:"key_id"`
	UserID       string `xorm:"user_id" json:"user_id"`
	EncryptedKey []byte `xorm:"encrypted_key" json:"encrypted_key"`
}

type SymmetricKeyOperation int

const (
	OperationCreateKey SymmetricKeyOperation = 1
	OperationGetKey    SymmetricKeyOperation = 2
	OperationRemoveKey SymmetricKeyOperation = 4
)

type LogSymmetricKey struct {
	KeyID     string                `xorm:"key_id" json:"key_id"`
	UserID    string                `xorm:"user_id" json:"user_id"`
	WhoId     string                `xorm:"who_id" json:"who_id"`
	Operation SymmetricKeyOperation `xorm:"operation" json:"operation"`
	Created   int64                 `xorm:"created" json:"created"`
}

type KeyUserPair struct {
	KeyID  string `xorm:"key_id" json:"key_id"`
	UserID string `xorm:"user_id" json:"user_id"`
}
