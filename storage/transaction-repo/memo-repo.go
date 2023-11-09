package transactionrepo

import (
	"context"
	"errors"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/types"
)

type MemoTransactionRepo struct {
	transactions []*types.Transaction
}

func NewMemoTransactionRepo() *MemoTransactionRepo {
	var transaction []*types.Transaction
	return &MemoTransactionRepo{
		transactions: transaction,
	}
}

func (tr *MemoTransactionRepo) CreateTransaction(ctx context.Context, transaction *types.Transaction) error {
	tr.transactions = append(tr.transactions, transaction)
	return nil
}

func (tr *MemoTransactionRepo) GetTransactionsHistory(ctx context.Context, accountId string, filters requestparams.GetTransactionsHistoryRequest) ([]*types.Transaction, error) {
	return tr.transactions, nil
}

func (tr *MemoTransactionRepo) GetTransaction(ctx context.Context, id string) (*types.Transaction, error) {
	for i, t := range tr.transactions {
		if t.ID == id {
			return tr.transactions[i], nil
		}
	}
	return nil, errors.New("")
}

func (tr *MemoTransactionRepo) GetMultiBeneficiaryTransactions(ctx context.Context, multibeneficiaryid string) ([]*types.Transaction, error) {
	return tr.transactions, nil
}
