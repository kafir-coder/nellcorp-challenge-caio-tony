package requestparams

type RefundMoneyRequest struct {
	IsMultibenificiary bool   `json:"is_multi_benificiary"`
	MultiBeneficiaryId string `json:"multi_beneficiary_id"`
	TransactionId      string `json:"transaction_id"`
}

func (req *RefundMoneyRequest) Validate() bool {
	if req.IsMultibenificiary && req.MultiBeneficiaryId == "" {
		return false
	}
	if !req.IsMultibenificiary && req.TransactionId == "" {
		return false
	}
	return true
}
