package services

import (
	"context"
	"database/sql"
	"errors"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"
	"go-sample/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccountSrv(t *testing.T) {

	accountRepo := accountrepo.NewMemoAccountRepo()
	transactionRepo := transactionrepo.NewMemoTransactionRepo()
	accountSrv := NewAccount(accountRepo, transactionRepo)
	t.Run("accounts.CreateAccount should call AccountRepo.CreateAccount", func(t *testing.T) {

		ctx := context.Background()
		account := &types.Account{}

		accountRepo.On("CreateAccount", ctx, account).Times(1).Return(nil)

		accountSrv.CreateAccount(ctx, account)
		accountRepo.AssertExpectations(t)

	})

	t.Run("srv.CreateAccount should return error if repo.CreateAccount returns error", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{}

		want := errors.New("some dumb error")
		accountRepo.On("CreateAccount", ctx, account).Return(want)

		have := accountSrv.CreateAccount(ctx, account)

		assert.True(t, errors.Is(want, have))
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.IsThereAlreadyAccountWithThisOwner should call AccountRepo.GetAccountByOwnerId", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{}

		accountRepo.On("GetAccountByOwnerId", ctx, account.Owner).Return(account, nil).Times(1)

		accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, account.ID)
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.IsThereAlreadyAccountThisOwner should return false if accountRepo.GetAccountByOwner returns a ErrNoRows error", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{}

		accountRepo.On("GetAccountByOwnerId", ctx, account.Owner).Times(1).Return(account, sql.ErrNoRows)
		have, _ := accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, account.ID)

		assert.False(t, have)
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.IsThereAlreadyAccountWithThisOwner return true if accountRepo.GetAccountByOwner returns nil error object", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{
			ID:        "account",
			Owner:     "account",
			Balance:   23,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		}

		accountRepo.On("GetAccountByOwnerId", ctx, account.Owner).Return(account, nil).Times(1)
		have, _ := accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, account.ID)

		assert.True(t, have)
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.GetAccountBalance should return InexistentAccount Error if accountRepo.GetAccountBalance returns ErrNoRows", func(t *testing.T) {
		ctx := context.Background()
		balance := 12
		accountId := ""

		accountRepo.On("GetAccountBalance", ctx, accountId).Return(balance, sql.ErrNoRows).Times(1)
		_, err := accountSrv.GetAccountBalance(ctx, accountId)

		assert.True(t, err == ErrInexistentAccount)

		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.GetAccountBalance should call repo.GetAccountBalance", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{
			ID:        "account",
			Owner:     "account",
			Balance:   23,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		}

		accountRepo.On("GetAccountBalance", ctx, account.Owner).Return(account.Balance, nil).Times(1)

		accountSrv.GetAccountBalance(ctx, account.ID)
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.GetAccountBalance should return error if repo.GetAccountBalance returns error", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{
			ID:        "account",
			Owner:     "account",
			Balance:   23,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		}

		want := errors.New("some dumb error")
		accountRepo.On("GetAccountBalance", ctx, account.Owner).Return(0, want).Times(1)

		_, have := accountSrv.GetAccountBalance(ctx, account.ID)

		assert.True(t, errors.Is(want, have))
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.GetAccountBalance should return the same balance as repo.AccountBalance Returns", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{
			ID:        "account",
			Owner:     "account",
			Balance:   23,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		}

		accountRepo.On("GetAccountBalance", ctx, account.Owner).Return(account.Balance, nil).Times(1)

		have, _ := accountSrv.GetAccountBalance(ctx, account.ID)

		assert.Equal(t, account.Balance, have)
		accountRepo.AssertExpectations(t)
	})

	t.Run("accounts.DepositMoney should call repo.IncrBalance", func(t *testing.T) {

		ctx := context.Background()
		accountId := "some dumb id"
		amount := 100
		accountRepo.On("IncrBalance", ctx, accountId, float64(amount)).Times(1).Return(nil)
		accountSrv.DepositMoney(ctx, accountId, float64(amount))

		accountRepo.AssertExpectations(t)
	})

}
