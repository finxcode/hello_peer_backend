package services

type NoPetFoundError struct {
	Msg string
}

func NewNoPetError(msg string) *NoPetFoundError {
	return &NoPetFoundError{
		Msg: msg,
	}
}

func (n *NoPetFoundError) Error() string {
	return n.Msg
}
