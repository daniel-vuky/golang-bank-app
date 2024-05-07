package db

import (
	"context"
	"testing"

	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/stretchr/testify/require"
)

func TestStore_TransferTx(t *testing.T) {

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	amount := util.RandomMoney()

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccountTx := result.FromAccount
		require.NotEmpty(t, fromAccountTx)
		require.Equal(t, fromAccount.ID, fromAccountTx.ID)

		toAccountTx := result.ToAccount
		require.NotEmpty(t, toAccountTx)
		require.Equal(t, toAccount.ID, toAccountTx.ID)

		// check balances
		fromAccountBalanceDiff := fromAccount.Balance - fromAccountTx.Balance //100 - 90
		toAccountBalanceDiff := toAccountTx.Balance - toAccount.Balance       // 110 - 100
		require.Equal(t, fromAccountBalanceDiff, toAccountBalanceDiff)
		require.True(t, fromAccountBalanceDiff > 0)
		require.True(t, fromAccountBalanceDiff%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(fromAccountBalanceDiff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	fromAccountUpdated, fromAccountUpdatedErr := testStore.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, fromAccountUpdatedErr)
	require.NotEmpty(t, fromAccountUpdated)

	toAccountUpdated, toAccountUpdatedErr := testStore.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, toAccountUpdatedErr)
	require.NotEmpty(t, toAccountUpdated)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, fromAccountUpdated.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, toAccountUpdated.Balance)
}
