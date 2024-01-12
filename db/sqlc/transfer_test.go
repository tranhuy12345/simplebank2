package db

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	argument := CreateTransfersParams{
		FromAccountID: sql.NullInt64{Int64: 1, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: 2, Valid: true},
		Amount:        1000,
		CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
	}

	transfer, err := testQueries.CreateTransfers(context.Background(), argument)
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
}
