package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestUtil_LoadConfig tests load config
func TestUtil_LoadConfig(t *testing.T) {
	config, err := LoadConfig("../")
	require.NoError(t, err)
	require.NotEmpty(t, config)
}
