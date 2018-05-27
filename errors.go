package toy_chain

import "errors"

var (
	ErrChainReplaceTooShort = errors.New("blockchain: replacement chain is too short")
	ErrChainReplaceInvalid  = errors.New("blockchain: replacement chain is invalid")

	ErrTransactionAmountExceedsBalance = errors.New("transaction: amount exceeds balance")
)
