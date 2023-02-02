package data

type LastBlock interface {
	Set(number int64, chain int64) error
	Get(chain int64) (*int64, error)
}
