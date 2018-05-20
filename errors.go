package toy_chain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrChainReplaceTooShort = Error("blockchain: replacement chain is too short")
	ErrChainReplaceInvalid = Error("blockchain: replacement chain is invalid")
)
