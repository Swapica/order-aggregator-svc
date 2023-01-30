-- +migrate Up
CREATE TABLE orders
(
    id            bigint,
    src_chain     bigint,
    account       varchar(42) NOT NULL,
    sell_token    varchar(42) NOT NULL,
    buy_token     varchar(42) NOT NULL,
    sell_amount   text        NOT NULL,
    buy_amount    text        NOT NULL,
    dest_chain    bigint      NOT NULL,

    state         smallint    NOT NULL,
    executed_by   bigint,
    match_swapica varchar(42),
    PRIMARY KEY (id, src_chain)
);

CREATE TABLE match_orders
(
    id          bigint,
    src_chain   bigint,
    order_id    bigint      NOT NULL,
    order_chain bigint      NOT NULL,
    account     varchar(42) NOT NULL,
    sell_token  varchar(42) NOT NULL,
    sell_amount text        NOT NULL,

    state       smallint    NOT NULL,
    PRIMARY KEY (id, src_chain),
    FOREIGN KEY (order_id, order_chain) REFERENCES orders (id, src_chain)
);

CREATE TABLE last_blocks
(
    number    bigint NOT NULL,
    src_chain bigint,
    CONSTRAINT last_blocks_pkey PRIMARY KEY (src_chain)
);

-- +migrate Down
DROP TABLE match_orders;
DROP TABLE orders;
DROP TABLE last_blocks;
