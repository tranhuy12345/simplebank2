package db

import (
	"context"
	"db/db/util"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func createdAccountRandom(t *testing.T) Accounts {
	user := createdUserRandom(t)
	argument := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), argument)
	if err != nil {
		log.Fatal(err)
	}
	require.Equal(t, argument.Owner, account.Owner)
	require.Equal(t, argument.Balance, account.Balance)
	require.Equal(t, argument.Currency, account.Currency)
	require.NotZero(t, account.ID)
	return account
}
func TestAccount(t *testing.T) {
	createdAccountRandom(t)
}

func TestListAccounts(t *testing.T) {
	var listAccounts []ListQueryRow
	listAccounts, err := testQueries.ListQuery(context.Background())
	for i := 0; i < 5; i++ {
		fmt.Println(listAccounts[i])
	}
	require.NoError(t, err)

}
