package accountrepo

import (
	"context"
	"go-sample/types"

	"github.com/jmoiron/sqlx"
)

type IAccountRepo interface {
	CreateAccount(context.Context, *types.Account) error
	GetAccountById(context.Context, string) (*types.Account, error)
	GetAccountByOwnerId(context.Context, string) (*types.Account, error)
	GetAccountBalance(context.Context, string) (float64, error)
	IncrBalance(context.Context, string, float64) error
	DecrBalance(context.Context, string, float64) error
}
type AccountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) AccountRepo {
	return AccountRepo{db}
}

func (ar AccountRepo) CreateAccount(ctx context.Context, account *types.Account) error {
	_, err := ar.db.ExecContext(ctx, "INSERT INTO accounts VALUES($1, $2, $3, $4, $5, $6)", account.ID, account.Owner, account.Balance, account.CreatedAt, account.UpdatedAt, account.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (ar AccountRepo) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	var account types.Account
	err := ar.db.QueryRowContext(ctx, "SELECT * FROM accounts WHERE id=$1", id).Scan(
		&account.ID, &account.Owner, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)

	return &account, err
}

func (ar AccountRepo) GetAccountByOwnerId(ctx context.Context, ownerId string) (*types.Account, error) {
	var account types.Account
	err := ar.db.QueryRowContext(ctx, "SELECT id, owner_id, balance::decimal, created_at, updated_at, deletedat FROM accounts WHERE owner_id = $1", ownerId).Scan(
		&account.ID, &account.Owner, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)

	return &account, err
}

func (ar AccountRepo) GetAccountBalance(ctx context.Context, accountId string) (float64, error) {
	var balance float64
	err := ar.db.QueryRowContext(ctx, "SELECT balance::decimal FROM accounts WHERE accounts.id = $1", accountId).Scan(&balance)
	return balance, err
}

func (ar AccountRepo) IncrBalance(ctx context.Context, accountId string, incr float64) error {
	_, err := ar.db.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", incr, accountId)
	return err
}

func (ar AccountRepo) DecrBalance(ctx context.Context, accountId string, incr float64) error {
	_, err := ar.db.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", incr, accountId)
	return err
}
