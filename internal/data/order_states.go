package data

const (
	StateAwaitingMatch uint8 = iota + 1
	StateAwaitingFinalization
	StateCanceled
	StateExecuted
	StateBadToken = 255
)

func StateToString(state uint8) string {
	switch state {
	case StateAwaitingMatch:
		return "Awaiting Match"
	case StateAwaitingFinalization:
		return "Awaiting Finalization"
	case StateCanceled:
		return "Canceled"
	case StateExecuted:
		return "Executed"
	case StateBadToken:
		return "Bad Token"
	default:
		return ""
	}
}
