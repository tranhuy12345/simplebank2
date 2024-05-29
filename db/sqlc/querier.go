// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountsBalance(ctx context.Context, arg AddAccountsBalanceParams) (Accounts, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Accounts, error)
	CreateEntries(ctx context.Context, arg CreateEntriesParams) (Entries, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Sessions, error)
	CreateTransfers(ctx context.Context, arg CreateTransfersParams) (Transfers, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (Users, error)
	DeleteAccounts(ctx context.Context, id int64) error
	DeleteEntries(ctx context.Context, id int64) error
	DeleteEntriesByAccountId(ctx context.Context, accountID sql.NullInt64) error
	DeleteTransfers(ctx context.Context, id int64) error
	DeleteTransfersByAccountId(ctx context.Context, accountID sql.NullInt64) error
	GetAccount(ctx context.Context, id int64) (Accounts, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Accounts, error)
	GetEntries(ctx context.Context, id int64) (Entries, error)
	GetSession(ctx context.Context, id uuid.UUID) (Sessions, error)
	GetTransfers(ctx context.Context, id int64) (Transfers, error)
	GetUser(ctx context.Context, username string) (Users, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Accounts, error)
	ListEntries(ctx context.Context) ([]Entries, error)
	ListQuery(ctx context.Context) ([]ListQueryRow, error)
	ListTransfers(ctx context.Context) ([]Transfers, error)
	UpdateAccounts(ctx context.Context, arg UpdateAccountsParams) (Accounts, error)
	UpdateEntries(ctx context.Context, arg UpdateEntriesParams) error
	UpdateTransfers(ctx context.Context, arg UpdateTransfersParams) error
}

var _ Querier = (*Queries)(nil)
