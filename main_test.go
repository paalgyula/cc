package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractIPFromBinary(t *testing.T) {
	name := "boob-Fa0224011388"

	c, err := resolveRemoteFromName(name)
	assert.NoError(t, err)
	assert.Equal(t, "250.2.36.1:5000", c)
}
