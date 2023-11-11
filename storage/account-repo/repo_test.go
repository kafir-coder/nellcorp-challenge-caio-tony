package accountrepo

import (
	"context"
	"database/sql"
	"fmt"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/types"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func createAccountTable(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.accounts (
		id VARCHAR(36) PRIMARY KEY,
		owner_id VARCHAR(36) NOT NULL,
		balance MONEY NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deletedAt TIMESTAMP NULL
	  );`)
	if err != nil {
		panic(err)
	}
}

func dropAccountTable(db *sqlx.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS public.accounts;`)
	if err != nil {
		panic(err)
	}
}
func TestAccountRepo(t *testing.T) {

	TEST_DATABASE_NAME := "nell_challenge_test"

	TEST_DATABASE_URL := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", TEST_DATABASE_NAME)

	db, _ := sqlx.Connect("postgres", TEST_DATABASE_URL)

	createAccountTable(db)

	repo := &AccountRepo{db}
	t.Run("Create Account", func(t *testing.T) {
		ctx := context.Background()
		owner_id := "some dumb id"
		balance := 100.00

		entity := types.NewAccount(owner_id, balance)

		accountId := entity.ID
		repo.CreateAccount(ctx, entity)

		var account types.Account
		db.QueryRowContext(ctx, "SELECT id, owner_id, balance::decimal, created_at, updated_at, deletedat FROM accounts WHERE id=$1", accountId).Scan(
			&account.ID, &account.Owner, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)

		if account.ID != entity.ID {
			t.Errorf("account.ID = %s, expected %s", account.ID, entity.ID)
		}
		t.Cleanup(func() {
			db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
		})
	})

	t.Run("GetAcountById", func(t *testing.T) {
		t.Run("should return an account", func(t *testing.T) {
			ctx := context.Background()
			owner_id := "some dumb id"
			balance := 100.00

			entity := types.NewAccount(owner_id, balance)

			accountId := entity.ID
			repo.CreateAccount(ctx, entity)

			acc, _ := repo.GetAccountById(ctx, accountId)

			if acc.ID != entity.ID {
				t.Errorf("acc.ID = %s, expected %s", acc.ID, entity.ID)
			}

			t.Cleanup(func() {
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
			})
		})
		t.Run("should return a sql.ErrNoRows", func(t *testing.T) {

			ctx := context.Background()
			accountId := "f57d2e6f-8a73-470b-9067-729cab4cabc5"

			_, err := repo.GetAccountById(ctx, accountId)

			if err != sql.ErrNoRows {
				t.Errorf("err = %s, expected %s", err, sql.ErrNoRows)
			}

		})
	})

	t.Run("GetAccoutBalance", func(t *testing.T) {
		t.Run("should return the account balance", func(t *testing.T) {
			ctx := context.Background()
			owner_id := "some dumb id"
			balance := 100.00

			entity := types.NewAccount(owner_id, balance)

			accountId := entity.ID
			repo.CreateAccount(ctx, entity)

			res, _ := repo.GetAccountBalance(ctx, accountId)

			if res != balance {
				t.Errorf("res = %f, expected %f", res, balance)
			}

			t.Cleanup(func() {
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
			})
		})

		t.Run("should return a sql.ErrNoRows", func(t *testing.T) {

			ctx := context.Background()
			accountId := "f57d2e6f-8a73-470b-9067-729cab4cabc5"

			_, err := repo.GetAccountBalance(ctx, accountId)

			if err != sql.ErrNoRows {
				t.Errorf("err = %s, expected %s", err, sql.ErrNoRows)
			}

		})
	})

	t.Run("GetAccountByOwnerId", func(t *testing.T) {
		t.Run("should return the account by owner id", func(t *testing.T) {
			ctx := context.Background()
			owner_id := "some dumb id"
			balance := 100.00

			entity := types.NewAccount(owner_id, balance)

			accountId := entity.ID
			repo.CreateAccount(ctx, entity)

			res, _ := repo.GetAccountByOwnerId(ctx, owner_id)

			if res.Owner != entity.Owner {
				t.Errorf("res.ID = %s, expected %s", res.ID, entity.ID)
			}

			t.Cleanup(func() {
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
			})
		})

		t.Run("should return a sql.ErrNoRows", func(t *testing.T) {

			ctx := context.Background()
			ownerId := "f57d2e6f-8a73-470b-9067-729cab4cabc5"

			_, err := repo.GetAccountByOwnerId(ctx, ownerId)

			if err != sql.ErrNoRows {
				t.Errorf("err = %s, expected %s", err, sql.ErrNoRows)
			}

		})
	})

	t.Run("ListAccounts", func(t *testing.T) {
		t.Run("should return the account by owner id", func(t *testing.T) {
			ctx := context.Background()

			accounts := []*types.Account{
				{
					ID:        "",
					Owner:     "1e55e90d-427e-4f0b-b71f-5e38e57c2ae8",
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: time.Now(),
				},
				{
					ID:        "",
					Owner:     "1e55e901-427e-4f0b-b71f-5e38e57c2ae8",
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: time.Now(),
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)

				accounts[i].ID = entity.ID
				repo.CreateAccount(ctx, entity)
			}

			request := requestparams.ListAccountsRequest{
				Page:  "1",
				Limit: "10",
			}
			res, _ := repo.ListAccounts(ctx, request)
			if len(res) != len(accounts) {
				t.Errorf("len(res.Accounts) = %d, expected %d", len(res), len(accounts))
			}
			t.Cleanup(func() {
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accounts[0].ID)
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accounts[1].ID)
			})
		})

		t.Run("should obey the limit filter", func(t *testing.T) {
			ctx := context.Background()

			accounts := []*types.Account{
				{
					ID:        "",
					Owner:     "1e55e90d-427e-4f0b-b71f-5e38e57c2ae8",
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: time.Now(),
				},
				{
					ID:        "",
					Owner:     "1e55e901-427e-4f0b-b71f-5e38e57c2ae8",
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: time.Now(),
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)
				accounts[i].ID = entity.ID
				repo.CreateAccount(ctx, entity)
			}

			limit := "1"

			request := requestparams.ListAccountsRequest{
				Page:  "1",
				Limit: limit,
			}
			res, _ := repo.ListAccounts(ctx, request)

			if limit, _ := strconv.Atoi(limit); len(res) != limit {
				t.Errorf("len(res.Accounts) = %d, expected %d", len(res), len(accounts))
			}

			t.Cleanup(func() {
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accounts[0].ID)
				db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accounts[1].ID)
			})
		})

		t.Run("should return an empty array", func(t *testing.T) {

			ctx := context.Background()

			request := requestparams.ListAccountsRequest{
				Page:  "1",
				Limit: "10",
			}
			accounts, err := repo.ListAccounts(ctx, request)
			if len(accounts) != 0 {
				t.Errorf("err = %s, expected %s", err, sql.ErrNoRows)
			}

		})
	})
	t.Run("IncrBalance", func(t *testing.T) {

		ctx := context.Background()
		owner_id := "some dumb id"
		balance := 100.00

		entity := types.NewAccount(owner_id, balance)

		accountId := entity.ID
		repo.CreateAccount(ctx, entity)

		incr := 140
		repo.IncrBalance(ctx, accountId, float64(incr))

		var finalBalance float64
		db.QueryRowContext(ctx, "SELECT balance::decimal FROM accounts WHERE id=$1", accountId).Scan(
			&finalBalance)

		if finalBalance != balance+float64(incr) {
			t.Errorf("finalBalance = %f, expected %f", finalBalance, balance+float64(incr))
		}

		t.Cleanup(func() {
			db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
		})
	})

	t.Run("DecrBalance", func(t *testing.T) {

		ctx := context.Background()
		owner_id := "some dumb id"
		balance := 100.00

		entity := types.NewAccount(owner_id, balance)

		accountId := entity.ID
		repo.CreateAccount(ctx, entity)

		decr := 10
		repo.DecrBalance(ctx, accountId, float64(decr))

		var finalBalance float64
		db.QueryRowContext(ctx, "SELECT balance::decimal FROM accounts WHERE id=$1", accountId).Scan(
			&finalBalance)

		if finalBalance != balance-float64(decr) {
			t.Errorf("finalBalance = %f, expected %f", finalBalance, balance-float64(decr))
		}

		t.Cleanup(func() {
			db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", accountId)
		})
	})

	t.Cleanup(func() {
		dropAccountTable(db)
	})
}
