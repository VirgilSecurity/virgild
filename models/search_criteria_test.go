package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ResolveScope_PassGlobal_ReturnGlobal(t *testing.T) {
	actual := ResolveScope("global")
	assert.Equal(t, "global", actual)
}

func Test_ResolveScope_PassApplication_ReturnApplication(t *testing.T) {
	actual := ResolveScope("application")
	assert.Equal(t, "application", actual)
}

func Test_ResolveScope_PassOther_ReturnApplication(t *testing.T) {
	actual := ResolveScope("test")
	assert.Equal(t, "application", actual)
}
