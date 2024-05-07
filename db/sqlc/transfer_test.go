package db

import (
	"context"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testStore.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	require.NotZerof(t, transfer.ID, "transfer id should not be zero")
	require.NotZerof(t, transfer.CreatedAt, "transfer created_at should not be zero")

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestQueries_UpdateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount, toAccount)
	args := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: util.RandomMoney(),
	}
	err := testStore.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)

	updatedTransfer, updatedErr := testStore.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, updatedErr)
	require.NotEmpty(t, updatedTransfer)
	require.Equal(t, transfer.ID, updatedTransfer.ID)
	require.Equal(t, args.Amount, updatedTransfer.Amount)
	require.Equal(t, transfer.FromAccountID, updatedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, updatedTransfer.ToAccountID)
	require.Equal(t, transfer.CreatedAt, updatedTransfer.CreatedAt)
}

func TestQueries_GetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount, toAccount)
	foundTransfer, err := testStore.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, foundTransfer)
	require.Equal(t, transfer.ID, foundTransfer.ID)
	require.Equal(t, transfer.FromAccountID, foundTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, foundTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, foundTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, foundTransfer.CreatedAt)
}

func TestQueries_ListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
	}
	args := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        0,
	}
	transfer, err := testStore.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfer, 5)

	for _, transfer := range transfer {
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
	}
}
