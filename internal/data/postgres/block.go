package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const blockTable = "last_blocks"
const upsertBlockSuffix = "ON CONFLICT ON CONSTRAINT last_blocks_pkey DO UPDATE SET number = ?"

type block struct {
	db *pgdb.DB
}

func NewLastBlock(db *pgdb.DB) data.LastBlock {
	return block{db: db}
}

func (q block) Set(n uint64, chain string) error {
	stmt := squirrel.Insert(blockTable).Columns("number", "src_chain").Values(n, chain).Suffix(upsertBlockSuffix, n)
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to initialize or update last block")
}

func (q block) Get(chain string) (*uint64, error) {
	var result struct {
		Number uint64 `db:"number"`
	}
	stmt := squirrel.Select("number").From(blockTable).Where(squirrel.Eq{"src_chain": chain})

	if err := q.db.Get(&result, stmt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select last block")
	}

	return &result.Number, nil
}
