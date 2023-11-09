package accountrepo

import (
	"context"
	"errors"
	"go-sample/types"

	"github.com/stretchr/testify/mock"
)

type MemoAccountRepo struct {
	accounts []*types.Account
	mock.Mock
}

func NewMemoAccountRepo() *MemoAccountRepo {
	accounts := []*types.Account{}
	return &MemoAccountRepo{
		accounts: accounts,
	}
}

func (m *MemoAccountRepo) CreateAccount(ctx context.Context, account *types.Account) error {
	m.accounts = append(m.accounts, account)
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MemoAccountRepo) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	args := m.Called(ctx, id)

	for _, account := range m.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return args.Get(0).(*types.Account), args.Error(0)
}

func (m *MemoAccountRepo) GetAccountByOwnerId(ctx context.Context, ownerId string) (*types.Account, error) {
	args := m.Called(ctx, ownerId)
	return args.Get(0).(*types.Account), args.Error(1)
}

func (m *MemoAccountRepo) GetAccountBalance(ctx context.Context, accountId string) (float64, error) {

	args := m.Called(ctx, accountId)
	s, _ := args.Get(0).(float64)
	return s, args.Error(1)
}

func (m *MemoAccountRepo) IncrBalance(ctx context.Context, accountId string, incr float64) error {
	args := m.Called(ctx, accountId, incr)
	return args.Error(0)
}

func (m *MemoAccountRepo) DecrBalance(ctx context.Context, accountId string, incr float64) error {
	for i, a := range m.accounts {
		if a.ID == accountId {
			m.accounts[i].Balance -= incr
			return nil
		}
	}
	return errors.New("")
}
