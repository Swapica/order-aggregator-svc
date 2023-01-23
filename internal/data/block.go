package data

type LastBlock interface {
	Set(number uint64, chain string) error
	Get(chain string) (*uint64, error)
}
