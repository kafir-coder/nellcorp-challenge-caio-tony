package accountrepo

import (
	"context"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/types"
)

type MockCreateAccount struct {
	Called                     bool
	Calls                      int
	ExpectedReturnError        error
	ExpectedReturnErrorMessage string
	ExpectedReturn             string
}

type MockGetAccountByOwnerId struct {
	Called                     bool
	Calls                      int
	ExpectedReturnError        error
	ExpectedReturnErrorMessage string
	ExpectedReturn             *types.Account
}

type MockGetAccountBalance struct {
	Called                     bool
	Calls                      int
	ExpectedReturnError        error
	ExpectedReturnErrorMessage string
	ExpectedReturn             float64
}

type MockMemoAccountRepo struct {
	accounts             []*types.Account
	McreateAccount       MockCreateAccount
	MgetAccountBYOwnerId MockGetAccountByOwnerId
	MgetAccountBalance   MockGetAccountBalance
}

func NewMockMemoAccountRepo() *MockMemoAccountRepo {
	accounts := []*types.Account{}
	return &MockMemoAccountRepo{
		accounts: accounts,
	}
}

func (m *MockMemoAccountRepo) CreateAccount(ctx context.Context, account *types.Account) (string, error) {
	m.McreateAccount.Called = true
	m.McreateAccount.Calls++
	return m.McreateAccount.ExpectedReturn, m.McreateAccount.ExpectedReturnError
}

func (m *MockMemoAccountRepo) ListAccounts(context.Context, requestparams.ListAccountsRequest) ([]*types.Account, error) {
	return nil, nil
}

func (m *MockMemoAccountRepo) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	return nil, nil
}

func (m *MockMemoAccountRepo) GetAccountByOwnerId(ctx context.Context, ownerId string) (*types.Account, error) {
	m.MgetAccountBYOwnerId.Called = true
	m.MgetAccountBYOwnerId.Calls++
	return m.MgetAccountBYOwnerId.ExpectedReturn, m.MgetAccountBYOwnerId.ExpectedReturnError
}

func (m *MockMemoAccountRepo) GetAccountBalance(ctx context.Context, accountId string) (float64, error) {
	m.MgetAccountBalance.Called = true
	m.MgetAccountBalance.Calls++
	return m.MgetAccountBalance.ExpectedReturn, m.MgetAccountBalance.ExpectedReturnError
}

func (m *MockMemoAccountRepo) IncrBalance(ctx context.Context, accountId string, incr float64) error {
	return nil
}

func (m *MockMemoAccountRepo) DecrBalance(ctx context.Context, accountId string, incr float64) error {
	return nil
}
