package services

import (
	"context"
	"database/sql"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/api/utils"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"
	"go-sample/types"

	"github.com/google/uuid"
)

type Account struct {
	accountRepo     accountrepo.IAccountRepo
	transactionRepo transactionrepo.ITransactionRepo
}

func NewAccount(accountRepo accountrepo.IAccountRepo, transactionRepo transactionrepo.ITransactionRepo) Account {
	return Account{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (as *Account) CreateAccount(ctx context.Context, account *types.Account) (string, error) {
	return as.accountRepo.CreateAccount(ctx, account)
}

func (as *Account) ListAccounts(ctx context.Context, filters requestparams.ListAccountsRequest) ([]*types.Account, error) {
	return as.accountRepo.ListAccounts(ctx, filters)
}

func (as *Account) GetAccountBalance(ctx context.Context, accountId string) (float64, error) {
	balance, err := as.accountRepo.GetAccountBalance(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInexistentAccount
		}
		return 0, err
	}
	return balance, nil
}

func (as *Account) GetAccount(ctx context.Context, accountId string) (*types.Account, error) {
	account, err := as.accountRepo.GetAccountById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInexistentAccount
		}
		return nil, err
	}
	return account, nil
}

func (as *Account) DepositMoney(ctx context.Context, accountId string, amount float64) error {

	err := as.accountRepo.IncrBalance(ctx, accountId, amount)

	if err != nil {
		return err
	}
	transaction := types.NewTransaction(amount, "", "", "Deposito", string(utils.DEPOSIT), "", false, "")
	err = as.transactionRepo.CreateTransaction(ctx, transaction)
	return err
}

func (as *Account) IsAccountExistent(ctx context.Context, accountId string) bool {
	_, err := as.accountRepo.GetAccountById(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}

func (as *Account) WithdrawMoney(ctx context.Context, accountId string, amount float64) error {
	err := as.accountRepo.DecrBalance(ctx, accountId, amount)
	if err != nil {
		return err
	}
	transaction := types.NewTransaction(amount, "", "", "Levantamento", string(utils.WITHDRAW), "", false, "")
	err = as.transactionRepo.CreateTransaction(ctx, transaction)
	return err
}

func (as *Account) IsThereAlreadyAccountWithThisOwner(ctx context.Context, ownerId string) (bool, error) {
	_, err := as.accountRepo.GetAccountByOwnerId(ctx, ownerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (as *Account) HasInsufficientFunds(ctx context.Context, accountId string, amountToOut float64) bool {
	balance, _ := as.accountRepo.GetAccountBalance(ctx, accountId)

	return amountToOut > balance
}

func (as *Account) TransferMoney(ctx context.Context, transferParams requestparams.TransferMoneyRequest) error {
	var multiBeneficiaryTransactionId string = ""

	if len(transferParams.Repcipients) > 1 {
		multiBeneficiaryTransactionId = uuid.NewString()
	}

	for _, recipient := range transferParams.Repcipients {
		_, err := as.accountRepo.GetAccountById(ctx, recipient.AccountId)

		if err == sql.ErrNoRows {
			return ErrInexistentAccount
		}
	}
	for _, recipient := range transferParams.Repcipients {

		err := as.transactionRepo.MakeTransferTransaction(ctx, transferParams.From, recipient.AccountId, recipient.Amount)

		if err != nil {
			return err
		}

		transaction := types.NewTransaction(
			recipient.Amount,
			transferParams.From,
			recipient.AccountId,
			transferParams.Subject,
			string(utils.TRANSFER),
			multiBeneficiaryTransactionId,
			false,
			"",
		)
		as.transactionRepo.CreateTransaction(ctx, transaction)
	}
	return nil
}

func (as *Account) GetTransactionsHistory(ctx context.Context, accountId string, filters requestparams.GetTransactionsHistoryRequest) ([]*types.Transaction, error) {
	return as.transactionRepo.GetTransactionsHistory(ctx, accountId, filters)
}

func (as *Account) RefundMoneyNormalTransfer(ctx context.Context, transactionId string) error {

	transaction, err := as.transactionRepo.GetTransaction(ctx, transactionId)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInexistentTransaction
		}
		return err
	}
	if transaction.IsRefund {
		return ErrUnableToRefundARefund
	}

	balance, _ := as.accountRepo.GetAccountBalance(ctx, transaction.To)
	if transaction.Amount > balance {
		return ErrInsufficientFunds
	}

	err = as.transactionRepo.MakeTransferTransaction(ctx, transaction.To, transaction.From, transaction.Amount)

	if err != nil {
		return err
	}

	newTransaction := types.NewTransaction(
		transaction.Amount,
		transaction.From,
		transaction.To,
		transaction.Subject,
		string(utils.REFUND), "",
		true,
		transactionId)
	as.transactionRepo.CreateTransaction(ctx, newTransaction)

	return nil
}

func (as *Account) RefundMultibeneficiaryTransfer(ctx context.Context, multibeneficiaryId string) error {
	transactions, err := as.transactionRepo.GetMultiBeneficiaryTransactions(ctx, multibeneficiaryId)

	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		if err != nil {
			return nil
		}

		if transaction.IsRefund {
			return ErrUnableToRefundARefund
		}
		err = as.transactionRepo.MakeTransferTransaction(ctx, transaction.To, transaction.From, transaction.Amount)

		if err != nil {
			return err
		}

		newTransaction := types.NewTransaction(
			transaction.Amount,
			transaction.From,
			transaction.To,
			transaction.Subject,
			string(utils.REFUND), "",
			true,
			transaction.ID)
		as.transactionRepo.CreateTransaction(ctx, newTransaction)
	}

	return nil
}
