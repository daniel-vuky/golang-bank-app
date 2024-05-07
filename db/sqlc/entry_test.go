package db

import (
	"context"
	"testing"

	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/stretchr/testify/require"
)

// createRandomAccount creates a random account for testing
func createRandomEntry(t *testing.T, account *Account) Entry {
	newEntry := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testStore.CreateEntry(context.Background(), newEntry)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, newEntry.AccountID, entry.AccountID)
	require.Equal(t, newEntry.Amount, entry.Amount)
	require.NotZerof(t, entry.ID, "entry id should not be zero")
	require.NotZerof(t, entry.CreatedAt, "entry created_at should not be zero")

	return entry
}

// TestQueries_CreateEntry tests the CreateEntry method
func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, &account)
}

// TestQueries_UpdateEntry tests the UpdateEntry method
func TestQueries_UpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, &account)
	args := UpdateEntryParams{
		ID:     entry.ID,
		Amount: util.RandomMoney(),
	}
	err := testStore.UpdateEntry(context.Background(), args)
	require.NoError(t, err)

	updatedEntry, updatedErr := testStore.GetEntry(context.Background(), entry.ID)
	require.NoError(t, updatedErr)
	require.NotEmpty(t, updatedEntry)
	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, args.Amount, updatedEntry.Amount)
	require.Equal(t, entry.AccountID, updatedEntry.AccountID)
	require.Equal(t, entry.CreatedAt, updatedEntry.CreatedAt)
}

// TestQueries_GetEntry tests the GetEntry method
func TestQueries_GetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, &account)
	foundEntry, err := testStore.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, foundEntry)
	require.Equal(t, entry.ID, foundEntry.ID)
	require.Equal(t, entry.AccountID, foundEntry.AccountID)
	require.Equal(t, entry.Amount, foundEntry.Amount)
	require.Equal(t, entry.CreatedAt, foundEntry.CreatedAt)
}

// TestQueries_ListEntries tests the ListEntries method
func TestQueries_ListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, &account)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testStore.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}

// TestQueries_DeleteEntry tests the DeleteEntry method
func TestQueries_DeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, &account)
	err := testStore.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	deletedEntry, deletedErr := testStore.GetEntry(context.Background(), entry.ID)
	require.Error(t, deletedErr)
	require.Empty(t, deletedEntry)
}
