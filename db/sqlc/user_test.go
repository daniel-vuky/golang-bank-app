package db

import (
	"context"
	"testing"

	"github.com/daniel-vuky/golang-bank-app/util"

	"github.com/stretchr/testify/require"
)

// createRandomUser creates a random user for testing
func createRandomUser(t *testing.T) User {
	hashedPassword, hashedPasswordErr := util.HashPassword(util.RandomString(6))
	require.NoError(t, hashedPasswordErr)
	user := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	createdUser, err := testStore.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.HashedPassword, createdUser.HashedPassword)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.NotEmpty(t, createdUser.Username)
	require.NotZero(t, createdUser)

	return createdUser
}

// TestQueries_CreateUser tests the CreateUser method
func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestQueries_GetUser tests the GetUser method
func TestQueries_GetUser(t *testing.T) {
	user := createRandomUser(t)
	foundUser, err := testStore.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, foundUser)
	require.Equal(t, user.Username, foundUser.Username)
	require.Equal(t, user.HashedPassword, foundUser.HashedPassword)
	require.Equal(t, user.FullName, foundUser.FullName)
	require.Equal(t, user.Email, foundUser.Email)
	require.NotEmpty(t, foundUser.CreatedAt)
}
