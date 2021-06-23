package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDOB(t *testing.T) {
	age := GetDOB(1992, 10, 21)
	assert.Equal(t, "1992-10-21 00:00:00 +0000 UTC", age.String())
}