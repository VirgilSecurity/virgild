package sqlmodels

type CardSql struct {
	Id           string `xorm:"PK"`
	Identity     string `xorm:"Index"`
	IdentityType string
	Scope        string
	Card         string
}
