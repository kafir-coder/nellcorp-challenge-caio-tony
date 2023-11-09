package utils

// transaction operation codes
type OP string

const (
	DEPOSIT  OP = "DEPOSIT"
	WITHDRAW OP = "WITHDRAW"
	TRANSFER OP = "TRANSFER"
	REFUND   OP = "REFUND"
)
