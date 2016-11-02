// +build integration
package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MakeDatabase_IncorrectStringFormat_Panic(t *testing.T) {
	assert.Panics(t, func() { MakeDatabase("db") })
}

func Test_MakeDatabase_IncorrectConnection_Panic(t *testing.T) {
	assert.Panics(t, func() { MakeDatabase("sqlite3:/") })
}

func Test_MakeDatabase_DriverNotRegistred_Panic(t *testing.T) {
	assert.Panics(t, func() { MakeDatabase("undefined:") })
}
