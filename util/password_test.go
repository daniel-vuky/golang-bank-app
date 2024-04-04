package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestHashPassword tests the HashPassword function.
// It generates a random password, hashes it, and checks that
func TestUtil_HashPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.Error(t, err)
}
