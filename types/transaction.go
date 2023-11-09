package types

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                            string    `json:"id"`
	From                          string    `json:"from"`
	To                            string    `json:"to"`
	Subject                       string    `json:"subject"`
	Operation                     string    `json:"operation"`
	Amount                        float64   `json:"amount"`
	MultiBeneficiaryTransactionId string    `json:"multi_Beneficiary_transaction_id"`
	IsRefund                      bool      `json:"is_refunded"`
	RefundedTransactionId         string    `json:"refunded_transaction_id"`
	CreatedAt                     time.Time `json:"created_at"`
}

func ValidateTransaction(u *Account) bool { return true }

func NewTransaction(amount float64, from string, to string, subject string, operation string, MultiBeneficiaryTransactionId string, isRefund bool, RefundedTransactionId string) *Transaction {
	return &Transaction{
		ID:                            uuid.NewString(),
		From:                          from,
		To:                            to,
		Subject:                       subject,
		Operation:                     operation,
		MultiBeneficiaryTransactionId: MultiBeneficiaryTransactionId,
		IsRefund:                      isRefund,
		RefundedTransactionId:         RefundedTransactionId,
		Amount:                        amount,
		CreatedAt:                     time.Now(),
	}
}

func (tr *Transaction) ValidateTransaction() bool {
	return tr.Subject == "" && tr.Operation == ""
}
