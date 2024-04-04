package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestUtil_RandomInt test function return a random int64
func TestUtil_RandomInt(t *testing.T) {
	var randomInt int64

	randomInt = RandomInt(1, 5)
	require.NotEmpty(t, randomInt)
	require.True(t, randomInt >= 1)
	require.True(t, randomInt <= 5)
}

// TestUtil_RandomString test function return a random string
func TestUtil_RandomString(t *testing.T) {
	var randomString string

	randomInt := int(RandomInt(1, 5))
	randomString = RandomString(randomInt)
	require.NotEmpty(t, randomString)
	require.Equal(t, randomInt, len(randomString))
	require.NotRegexp(t, randomString, `\0-9`)
}

// TestUtil_RandomCurrency test function return a random currency string (USD or EUR)
func TestUtil_RandomCurrency(t *testing.T) {
	var randomCurrency string

	randomCurrency = RandomCurrency()
	require.NotEmpty(t, randomCurrency)
	require.Contains(t, []string{"USD", "EUR"}, randomCurrency)
}
