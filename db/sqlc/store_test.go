package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	store := NewStore(testDB)

	account1 := createdAccountRandom(t)
	account2 := createdAccountRandom(t)
	fmt.Println(">>Before: ", account1.Balance, account2.Balance)
	//Chạy n lần transaction
	n := 3
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d: ", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransfersTX(ctx, TransfersTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//Check Tranfer
		transfer := result.Transfers
		require.NotEmpty(t, transfer)
		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, transfer.FromAccountID)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfers(context.Background(), transfer.ID)
		require.NoError(t, err)
		//Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntries(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntries(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, account1.ID, fromAccount.ID)
		// require.Equal(t, account1.Balance-amount, fromAccount.Balance)
		fmt.Println(">>tx: ", fromAccount.Balance, toAccount.Balance)
		dif1 := account1.Balance - fromAccount.Balance
		fmt.Println(">>dif1: ", dif1)
		//dif2 := toAccount.Balance - account2.Balance
		//require.Equal(t, dif1, dif2)
		//require.True(t, dif1 > 0) //So tien chuyen phai lon hon 0
		//require.True(t, (dif1%amount) == 0)

		//k := int(dif1 / amount)
		//require.True(t, k >= 1 && k <= n)
	}

	updateAccounts1, err := store.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccounts2, err := store.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>After: ", updateAccounts1.Balance, updateAccounts2.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccounts2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccounts1.Balance)

}

func TestTransactionTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createdAccountRandom(t)
	account2 := createdAccountRandom(t)
	fmt.Println(">>Before: ", account1.Balance, account2.Balance)
	//Chạy n lần transaction
	n := 4
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d: ", i+1)
		fromAccount := account1.ID
		toAccount := account2.ID

		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransfersTX(ctx, TransfersTxParams{
				FromAccountID: fromAccount,
				ToAccountID:   toAccount,
				Amount:        amount,
			})

			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updateAccounts1, err := store.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccounts2, err := store.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>After: ", updateAccounts1.Balance, updateAccounts2.Balance)
	require.Equal(t, account2.Balance, updateAccounts2.Balance)
	require.Equal(t, account1.Balance, updateAccounts1.Balance)

}
