package services

import (
	"context"
	"database/sql"
	"errors"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"
	"go-sample/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountSrv(t *testing.T) {

	accountRepo := accountrepo.NewMockMemoAccountRepo()
	transactionRepo := transactionrepo.NewMemoTransactionRepo()
	accountSrv := NewAccount(accountRepo, transactionRepo)
	t.Run("accounts.CreateAccount should call AccountRepo.CreateAccount", func(t *testing.T) {

		ctx := context.Background()
		account := &types.Account{}
		got := accountRepo.McreateAccount.Calls
		expected := 1
		accountSrv.CreateAccount(ctx, account)

		if !accountRepo.McreateAccount.Called && expected != got {
			t.Errorf("Expected 1 call(s) but got %d calls", expected)
		}
	})
	t.Run("srv.CreateAccount should return error if repo.CreateAccount returns error", func(t *testing.T) {
		ctx := context.Background()
		account := &types.Account{}

		want := errors.New("some dumb error")

		accountRepo.McreateAccount.ExpectedReturnError = want

		_, have := accountSrv.CreateAccount(ctx, account)

		assert.True(t, errors.Is(want, have))
	})
	t.Run("accounts.IsThereAlreadyAccountWithThisOwner should call AccountRepo.GetAccountByOwnerId", func(t *testing.T) {
		ctx := context.Background()
		ownerId := "some dumb id"
		got := accountRepo.MgetAccountBYOwnerId.Calls
		expected := 1
		accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, ownerId)

		if !accountRepo.MgetAccountBYOwnerId.Called && expected != got {
			t.Errorf("Expected %d call(s) but got %d calls", expected, got)
		}
	})
	t.Run("accounts.IsThereAlreadyAccountThisOwner should return false if accountRepo.GetAccountByOwner returns a ErrNoRows error", func(t *testing.T) {
		ctx := context.Background()
		ownerId := "some dumb id"
		expected := false

		accountRepo.MgetAccountBYOwnerId.ExpectedReturnError = sql.ErrNoRows

		got, _ := accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, ownerId)

		if got != expected {
			t.Errorf("Expected %v but got %v calls", expected, got)
		}
	})
	t.Run("accounts.IsThereAlreadyAccountWithThisOwner return true if accountRepo.GetAccountByOwner returns nil error object", func(t *testing.T) {
		ctx := context.Background()
		ownerId := "some dumb id"
		expected := true

		accountRepo.MgetAccountBYOwnerId.ExpectedReturnError = nil

		got, _ := accountSrv.IsThereAlreadyAccountWithThisOwner(ctx, ownerId)

		if got != expected {
			t.Errorf("Expected %v but got %v calls", expected, got)
		}
	})
	t.Run("accounts.GetAccountBalance should call repo.GetAccountBalance", func(t *testing.T) {
		ctx := context.Background()
		accountId := "some dumb id"

		got := accountRepo.MgetAccountBalance.Calls
		expected := 1

		accountSrv.GetAccountBalance(ctx, accountId)

		if !accountRepo.MgetAccountBalance.Called && got != expected {
			t.Errorf("Expected %v but got %v calls", expected, got)
		}
	})
	t.Run("accounts.GetAccountBalance should return InexistentAccount Error if accountRepo.GetAccountBalance returns ErrNoRows", func(t *testing.T) {
		ctx := context.Background()
		accountId := "some dumb id"
		expected := ErrInexistentAccount

		accountRepo.MgetAccountBalance.ExpectedReturnError = sql.ErrNoRows

		_, got := accountSrv.GetAccountBalance(ctx, accountId)

		if got != expected {
			t.Errorf("Expected %v but got %v calls", expected, got)
		}
	})
	t.Run("accounts.GetAccountBalance should return the same balance as repo.AccountBalance Returns", func(t *testing.T) {
		ctx := context.Background()
		accountId := "some dumb id"
		expected := accountRepo.MgetAccountBalance.ExpectedReturn

		got, _ := accountSrv.GetAccountBalance(ctx, accountId)

		if got != expected {
			t.Errorf("Expected %v but got %v calls", expected, got)
		}
	})

}
