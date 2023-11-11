package transactionrepo

import (
	"context"
	"fmt"
	"go-sample/api/utils"
	accountrepo "go-sample/storage/account-repo"
	"go-sample/types"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func createTables(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.transaction_ (
		id VARCHAR(36) PRIMARY KEY,
		from_account VARCHAR(36) NOT NULL,
		to_account VARCHAR(36) NOT NULL,
		tr_status TEXT NOT NULL,
		operation VARCHAR(10) NOT NULL,
		amount MONEY NOT NULL,
		multiBeneficiaryId VARCHAR(36) NOT NULL,
		is_Refund BOOLEAN NOT NULL, 
		refunded_transaction_id VARCHAR(36) NOT NULL,
		createdAt TIMESTAMP NOT NULL
	  );`)

	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
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

func dropTables(db *sqlx.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS public.transaction_;`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS public.accounts;`)
	if err != nil {
		panic(err)
	}
}
func TestTransactionRepo(t *testing.T) {

	TEST_DATABASE_NAME := "nell_challenge_test"
	TEST_DATABASE_URL := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", TEST_DATABASE_NAME)
	db, _ := sqlx.Connect("postgres", TEST_DATABASE_URL)

	createTables(db)

	repo := &TransactionRepo{db}
	accountRepo := accountrepo.NewAccountRepo(db)

	t.Run("Create Account", func(t *testing.T) {
		ctx := context.Background()
		amount := 100.0
		subject := "some dumb subject"
		from := "some dumb from"
		to := "some dumb to"

		entity := types.NewTransaction(amount, from, to, subject, string(utils.DEPOSIT), "", false, "")

		transactionId := entity.ID
		repo.CreateTransaction(ctx, entity)

		var transaction types.Transaction
		db.QueryRowContext(
			ctx,
			`SELECT
				id, from_account, to_account, tr_status, operation, amount::decimal, multibeneficiaryid, is_refund, refunded_transaction_id, createdat
				FROM transaction_ WHERE id = $1`, entity.ID).Scan(
			&transaction.ID,
			&transaction.From,
			&transaction.To,
			&transaction.Subject,
			&transaction.Operation,
			&transaction.Amount,
			&transaction.MultiBeneficiaryTransactionId,
			&transaction.IsRefund,
			&transaction.RefundedTransactionId,
			&transaction.CreatedAt,
		)

		if transaction.ID != entity.ID {
			t.Errorf("account.ID = %s, expected %s", transaction.ID, entity.ID)
		}
		t.Cleanup(func() {
			db.ExecContext(ctx, "DELETE FROM transaction_ WHERE id=$1", transactionId)
		})
	})

	t.Run("MakeTransferTransaction", func(t *testing.T) {
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
				Balance:   0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: time.Now(),
			},
		}

		for i, account := range accounts {
			entity := types.NewAccount(account.Owner, account.Balance)

			accounts[i].ID = entity.ID
			accountRepo.CreateAccount(ctx, entity)
		}

		amount := 100.0
		from := accounts[0].ID
		to := accounts[1].ID

		repo.MakeTransferTransaction(ctx, from, to, amount)

		fromBalance, _ := accountRepo.GetAccountBalance(ctx, from)
		toBalance, _ := accountRepo.GetAccountBalance(ctx, to)

		if fromBalance != toBalance-amount {
			t.Errorf("fromBalance = %f, expected %f", fromBalance, amount)
		}
	})
	t.Cleanup(func() {
		dropTables(db)
	})
}
