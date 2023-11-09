package services

import "errors"

var ErrInexistentAccount = errors.New("InexistentAccount")
var ErrUnableToRefundARefund = errors.New("UnableToRefundARefundItself")
var ErrInsufficientFunds = errors.New("InsufficientFunds")
