package transactionrepo

import (
	"context"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/types"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type ITransactionRepo interface {
	CreateTransaction(context.Context, *types.Transaction) error
	GetTransactionsHistory(context.Context, string, requestparams.GetTransactionsHistoryRequest) ([]*types.Transaction, error)
	GetTransaction(context.Context, string) (*types.Transaction, error)
	GetMultiBeneficiaryTransactions(context.Context, string) ([]*types.Transaction, error)
}
type TransactionRepo struct {
	db *sqlx.DB
}

func NewATransactionRepo(db *sqlx.DB) TransactionRepo {
	return TransactionRepo{db}
}

func (tr TransactionRepo) CreateTransaction(ctx context.Context, transaction *types.Transaction) error {
	tx, err := tr.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO transaction_(id, from_account, to_account, tr_status, operation, amount, multibeneficiaryid, is_refund, refunded_transaction_id, createdat) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		transaction.ID,
		transaction.From,
		transaction.To,
		transaction.Subject,
		transaction.Operation,
		transaction.Amount,
		transaction.MultiBeneficiaryTransactionId,
		transaction.IsRefund,
		transaction.RefundedTransactionId,
		transaction.CreatedAt)

	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (tr TransactionRepo) GetTransactionsHistory(ctx context.Context, accountId string, filters requestparams.GetTransactionsHistoryRequest) ([]*types.Transaction, error) {

	page, _ := strconv.ParseInt(filters.Page, 10, 64)
	limit, _ := strconv.ParseInt(filters.Limit, 10, 64)

	offset := limit * (page - 1)
	var transactions []*types.Transaction
	rows, err := tr.db.QueryContext(
		ctx,
		`SELECT 
			id, from_account, to_account, tr_status, operation, amount::decimal, multibeneficiaryid, createdat FROM transaction_ 
			WHERE from_account = $1 LIMIT $2  OFFSET $3`, accountId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction types.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.From,
			&transaction.To,
			&transaction.Subject,
			&transaction.Operation,
			&transaction.Amount,
			&transaction.MultiBeneficiaryTransactionId,
			&transaction.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (tr TransactionRepo) GetTransaction(ctx context.Context, id string) (*types.Transaction, error) {
	var transaction types.Transaction
	err := tr.db.QueryRowContext(
		ctx,
		`SELECT 
			id, from_account, to_account, tr_status, operation, amount::decimal, multibeneficiaryid, is_refund, refunded_transaction_id, createdat 
			FROM transaction_ WHERE id = $1`, id).Scan(
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
	return &transaction, err
}

func (tr TransactionRepo) GetMultiBeneficiaryTransactions(ctx context.Context, multibeneficiaryid string) ([]*types.Transaction, error) {

	var transactions []*types.Transaction
	rows, err := tr.db.QueryContext(
		ctx, `SELECT 
				id, from_account, to_account, tr_status, operation, amount::decimal, multibeneficiaryid, is_refund, refunded_transaction_id, createdat 
				FROM transaction_ WHERE multibeneficiaryid = $1`, multibeneficiaryid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction types.Transaction
		err = rows.Scan(
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

		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}
