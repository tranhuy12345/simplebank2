package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEntries(t *testing.T) {

	var account Accounts

	accounted, err := testQueries.GetAccountForUpdate(context.Background(), 3)
	if err != nil {
		log.Fatal(err)
	}
	account = accounted
	require.Equal(t, int64(3), account.ID)
	argument := CreateEntriesParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Amount:    100,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	//log.Fatal(sql.NullInt64{Int64: account.ID})
	entries, err := testQueries.CreateEntries(context.Background(), argument)
	if err != nil {
		log.Fatal(err)
	}
	require.Equal(t, int64(3), entries.AccountID.Int64)
	fmt.Println(entries)
}

// func TestDeleteEntries(t *testing.T) {
// 	err := testQueries.DeleteAccounts(context.Background(), 3)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	log.Fatal(err)
// 	//require.Error(t, nil, err)
// }
