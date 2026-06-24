package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/paalgyula/cc/config"
)

func TestExtractIPFromBinary(t *testing.T) {
	name := "boob-Fa0224011388"

	c, err := resolveRemoteFromName(name)
	assert.NoError(t, err)
	assert.Equal(t, "250.2.36.1:5000", c)
}

func TestResolveRemoteFromName_ShortAddress(t *testing.T) {
	name := "foo-aabbccdd"

	_, err := resolveRemoteFromName(name)
	assert.ErrorContains(t, err, "must be 12 hex characters")
}

func TestResolveRemoteFromName_NoDash(t *testing.T) {
	_, err := resolveRemoteFromName("justaname")
	assert.ErrorIs(t, err, ErrNoAddrInFilename)
}

func TestResolveRemoteFromName_InvalidHex(t *testing.T) {
	name := "foo-ZZZZZZZZZZZZ"

	_, err := resolveRemoteFromName(name)
	assert.Error(t, err)
}

func TestFindShell_ReturnsNonEmpty(t *testing.T) {
	shell := findShell()
	assert.NotEmpty(t, shell)
}

func TestFindShell_UsesEnvVar(t *testing.T) {
	const expected = "/bin/sh"

	t.Setenv("SHELL", expected)

	shell := findShell()
	assert.Equal(t, expected, shell)
}

func TestFindShell_EnvVarNotExist(t *testing.T) {
	t.Setenv("SHELL", "/nonexistent/shell")

	shell := findShell()
	// Must fall back to a shell that exists
	assert.NotEmpty(t, shell)
	_, err := os.Stat(shell)
	assert.NoError(t, err, "shell %q must exist on disk", shell)
}

func TestResolveRemote_WithEnvVar(t *testing.T) {
	const expected = "1.2.3.4:8080"

	t.Setenv("R", expected)

	err := resolveRemote()
	require.NoError(t, err)
	assert.Equal(t, expected, config.Remote)
}

func TestResolveRemote_Empty(t *testing.T) {
	os.Clearenv()

	config.Remote = "" // ensure clean state
	err := resolveRemote()
	assert.ErrorIs(t, err, ErrNoRemote)
}

func TestResolveRemoteFromName_MinimalEdge(t *testing.T) {
	// Exactly the boundary case: binary name like "a-000000000000"
	c, err := resolveRemoteFromName("a-000000000000")
	require.NoError(t, err)
	assert.Equal(t, "0.0.0.0:0", c)
}
