package db

import (
	"context"
	"db/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createdUserRandom(t *testing.T) Users {
	hashedPassword, err := util.HashPassword(util.RandomString(10))
	argument := CreateUserParams{
		Username:     util.RandomOwner(),
		HashPassword: hashedPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), argument)
	require.NoError(t, err)
	require.Equal(t, argument.Username, user.Username)
	require.Equal(t, argument.HashPassword, user.HashPassword)
	require.Equal(t, argument.FullName, user.FullName)
	require.NotEmpty(t, user.Username)
	return user
}

func TestCreateUser(t *testing.T) {
	createdUserRandom(t)
}
