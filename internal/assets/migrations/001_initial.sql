-- +migrate Up
CREATE TABLE orders(
    id text,
    src_chain text,
    account varchar(42) NOT NULL,
    sell_tokens varchar(42) NOT NULL,
    buy_tokens varchar(42) NOT NULL,
    sell_amount text NOT NULL,
    buy_amount text NOT NULL,
    dest_chain text NOT NULL,

    state smallint NOT NULL,
    executed_by text,
    match_swapica text,
    PRIMARY KEY (id, src_chain)
);

CREATE TABLE last_blocks(
    number text NOT NULL,
    src_chain text,
    CONSTRAINT last_blocks_pkey PRIMARY KEY (src_chain)
);

-- +migrate Down
DROP TABLE orders;
DROP TABLE last_blocks;
