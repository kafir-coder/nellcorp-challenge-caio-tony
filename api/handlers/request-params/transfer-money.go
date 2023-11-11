package requestparams

type Recipient struct {
	AccountId string
	Amount    float64
}
type TransferMoneyRequest struct {
	From        string      `json:"from"`
	Amount      float64     `json:"amount"`
	Repcipients []Recipient `json:"recipients"`
	Subject     string      `json:"subject"`
}

func (t *TransferMoneyRequest) Validate() bool {

	var amountSum float64
	for _, recipient := range t.Repcipients {
		amountSum += recipient.Amount
	}

	if t.Amount <= 0 {
		return false
	}

	if amountSum != t.Amount {
		return false
	}
	return true
}
