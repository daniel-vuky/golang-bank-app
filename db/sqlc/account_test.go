package db

import (
	"context"
	"testing"

	"github.com/daniel-vuky/golang-bank-app/util"

	"github.com/stretchr/testify/require"
)

// createRandomAccount creates a random account for testing
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	account := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	createdAccount, err := testStore.CreateAccount(context.Background(), account)
	require.NoError(t, err)
	require.NotEmptyf(t, createdAccount, "created account should not be empty")
	require.Equal(t, account.Owner, createdAccount.Owner)
	require.Equal(t, account.Balance, createdAccount.Balance)
	require.Equal(t, account.Currency, createdAccount.Currency)
	require.NotZerof(t, createdAccount.ID, "created account id should not be zero")
	require.NotZerof(t, createdAccount.CreatedAt, "created account created_at should not be zero")

	return createdAccount
}

// TestQueries_CreateAccount tests the CreateAccount method
func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// TestQueries_UpdateAccount tests the UpdateAccount method
func TestQueries_UpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:       account.ID,
		Owner:    account.Owner,
		Balance:  util.RandomMoney(),
		Currency: account.Currency,
	}
	err := testStore.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	updatedAccount, updatedErr := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, updatedErr)
	require.NotEmptyf(t, updatedAccount, "updated account should not be empty")
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, args.Balance, updatedAccount.Balance)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, account.Currency, updatedAccount.Currency)
	require.Equal(t, account.CreatedAt, updatedAccount.CreatedAt)
}

func TestQueries_UpdateAccountBalance(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountBalanceParams{
		ID:     account.ID,
		Amount: util.RandomMoney(),
	}
	_, err := testStore.UpdateAccountBalance(context.Background(), args)
	require.NoError(t, err)

	updatedAccount, updatedErr := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, updatedErr)
	require.NotEmptyf(t, updatedAccount, "updated account should not be empty")
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, args.Amount+account.Balance, updatedAccount.Balance)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, account.Currency, updatedAccount.Currency)
	require.Equal(t, account.CreatedAt, updatedAccount.CreatedAt)
}

// TestQueries_GetAccount tests the GetAccount method
func TestQueries_GetAccount(t *testing.T) {
	account := createRandomAccount(t)
	foundAccount, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmptyf(t, foundAccount, "found account should not be empty")
	require.Equal(t, account.ID, foundAccount.ID)
	require.Equal(t, account.Balance, foundAccount.Balance)
	require.Equal(t, account.Owner, foundAccount.Owner)
	require.Equal(t, account.Currency, foundAccount.Currency)
	require.Equal(t, account.CreatedAt, foundAccount.CreatedAt)
}

// TestQueries_ListAccounts tests the ListAccounts method
func TestQueries_ListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	listAccountsParams := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), listAccountsParams)
	require.NoError(t, err)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

// TestQueries_DeleteAccount tests the DeleteAccount method
func TestQueries_DeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, deletedErr := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, deletedErr)
	require.Empty(t, deletedAccount)
}
