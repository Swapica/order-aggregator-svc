-- +migrate Up
CREATE TABLE orders
(
    id            text,
    src_chain     text,
    account       varchar(42) NOT NULL,
    sell_token    varchar(42) NOT NULL,
    buy_token     varchar(42) NOT NULL,
    sell_amount   text        NOT NULL,
    buy_amount    text        NOT NULL,
    dest_chain    text        NOT NULL,

    state         smallint    NOT NULL,
    executed_by   text,
    match_swapica text,
    PRIMARY KEY (id, src_chain)
);

CREATE TABLE match_orders
(
    id              text,
    src_chain       text,
    origin_order_id text        NOT NULL,
    account         varchar(42) NOT NULL,
    sell_token      varchar(42) NOT NULL,
    sell_amount     text        NOT NULL,
    origin_chain    text        NOT NULL,
    state           smallint    NOT NULL,

    PRIMARY KEY (id, src_chain),
    FOREIGN KEY (origin_order_id, src_chain) REFERENCES orders (id, src_chain)
);

CREATE TABLE last_blocks
(
    number    text NOT NULL,
    src_chain text,
    CONSTRAINT last_blocks_pkey PRIMARY KEY (src_chain)
);

-- +migrate Down
DROP TABLE match_orders;
DROP TABLE orders;
DROP TABLE last_blocks;
