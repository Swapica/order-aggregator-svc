package data

const (
	StateAwaitingMatch uint8 = iota + 1
	StateAwaitingFinalization
	StateCanceled
	StateExecuted
	StateBadToken = 255
)
