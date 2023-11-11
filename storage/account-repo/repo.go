package accountrepo

import (
	"context"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/types"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type IAccountRepo interface {
	CreateAccount(context.Context, *types.Account) (string, error)
	ListAccounts(context.Context, requestparams.ListAccountsRequest) ([]*types.Account, error)
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

func (ar AccountRepo) CreateAccount(ctx context.Context, account *types.Account) (string, error) {
	_, err := ar.db.ExecContext(ctx, "INSERT INTO accounts VALUES($1, $2, $3, $4, $5, $6)", account.ID, account.Owner, account.Balance, account.CreatedAt, account.UpdatedAt, account.DeletedAt)

	if err != nil {
		return "", err
	}
	var accountId string
	err = ar.db.QueryRowContext(ctx, "SELECT id FROM accounts WHERE id = $1", account.ID).Scan(&accountId)

	return accountId, err
}

func (ar AccountRepo) ListAccounts(ctx context.Context, filters requestparams.ListAccountsRequest) ([]*types.Account, error) {
	page, _ := strconv.ParseInt(filters.Page, 10, 64)
	limit, _ := strconv.ParseInt(filters.Limit, 10, 64)

	offset := limit * (page - 1)

	var accounts []*types.Account

	rows, err := ar.db.QueryContext(
		ctx,
		`SELECT 
			id, owner_id, balance::decimal, created_at, updated_at, deletedat FROM accounts 
			LIMIT $1  OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var account types.Account
		err = rows.Scan(
			&account.ID,
			&account.Owner,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt,
		)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (ar AccountRepo) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	var account types.Account
	err := ar.db.QueryRowContext(ctx, "SELECT id, owner_id, balance::decimal, created_at, updated_at, deletedat FROM accounts WHERE id=$1", id).Scan(
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
