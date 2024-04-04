package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestIsSupportCurrency tests the IsSupportCurrency function.
// It checks that it returns true for supported currencies, and false for unsupported currencies.
func TestUtil_IsSupportCurrency(t *testing.T) {
	var isSupportCurrency bool

	isSupportCurrency = IsSupportCurrency("USD")
	require.True(t, isSupportCurrency)

	isSupportCurrency = IsSupportCurrency("EUR")
	require.True(t, isSupportCurrency)

	isSupportCurrency = IsSupportCurrency("JPY")
	require.False(t, isSupportCurrency)
}
