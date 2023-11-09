package requestparams

type RefundMoneyRequest struct {
	IsMultibenificiary bool   `json:"is_multi_benificiary"`
	MultiBeneficiaryId string `json:"multi_beneficiary_id"`
	TransactionId      string `json:"transaction_id"`
}
