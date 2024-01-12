package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	*Queries         // Kế thừa những pt từ Queries
	db       *sql.DB //dùng để gọi Transaction
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// fnc : 1 chuỗi các queries
func (store *Store) execTx(ctx context.Context, fnc func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := New(tx)
	err = fnc(queries)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransfersTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfers   Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransfersTX(ctx context.Context, arg TransfersTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	//Bat dau transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		//Tao giao dich
		fmt.Println(txName, "create transfer")
		result.Transfers, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		//Luu entry cua nguoi chuyen
		fmt.Println(txName, "create entries 1")
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			Amount:    (-arg.Amount),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
		//Luu entry cua nguoi nhan
		fmt.Println(txName, "create entries 2")
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:    (arg.Amount),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txName, " arg.FromAccountID < arg.ToAccountID :Update account 1")
			result.FromAccount, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
				Amount: (-arg.Amount),
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}

			fmt.Println(txName, "arg.FromAccountID < arg.ToAccountID: Update account 2")
			result.ToAccount, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return nil
			}
		} else {
			fmt.Println(txName, "arg.FromAccountID > arg.ToAccountID: Update account 2")
			result.ToAccount, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return nil
			}
			fmt.Println(txName, "arg.FromAccountID > arg.ToAccountID: Update account 1")
			result.FromAccount, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
				Amount: (-arg.Amount),
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
