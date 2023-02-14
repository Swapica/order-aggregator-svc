-- +migrate Up
CREATE TABLE orders
(
    id            bigserial PRIMARY KEY,
    order_id      bigint      NOT NULL,
    src_chain     bigint      NOT NULL,
    creator       varchar(42) NOT NULL,
    sell_token    varchar(42) NOT NULL,
    buy_token     varchar(42) NOT NULL,
    sell_amount   numeric(78) NOT NULL,
    buy_amount    numeric(78) NOT NULL,
    dest_chain    bigint      NOT NULL,
    state         smallint    NOT NULL,
    match_id   bigint,
    match_swapica varchar(42),

    UNIQUE (order_id, src_chain)
);

CREATE TABLE match_orders
(
    id          bigserial PRIMARY KEY,
    match_id    bigint      NOT NULL,
    src_chain   bigint      NOT NULL,
    order_id    bigint      NOT NULL,
    order_chain bigint      NOT NULL,
    creator     varchar(42) NOT NULL,
    sell_token  varchar(42) NOT NULL,
    sell_amount numeric(78) NOT NULL,
    state       smallint    NOT NULL,

    UNIQUE (match_id, src_chain),
    FOREIGN KEY (order_id, order_chain) REFERENCES orders (order_id, src_chain)
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
