-- +migrate Up
CREATE TABLE tokens
(
    id        bigserial PRIMARY KEY,
    address   varchar(42) NOT NULL,
    src_chain bigint      NOT NULL,
    name      text        NOT NULL,
    symbol    text        NOT NULL,
    decimals  smallint    NOT NULL,

    UNIQUE (address, src_chain)
);

CREATE TABLE orders
(
    id                bigserial PRIMARY KEY,
    order_id          bigint      NOT NULL,
    src_chain         bigint      NOT NULL,
    creator           varchar(42) NOT NULL,
    sell_token        bigint      NOT NULL,
    buy_token         bigint      NOT NULL,
    sell_amount       numeric(78) NOT NULL,
    buy_amount        numeric(78) NOT NULL,
    dest_chain        bigint      NOT NULL,
    state             smallint    NOT NULL,
    use_relayer       boolean     NOT NULL,
    executed_by_match bigint,
    match_id          bigint,
    match_swapica     varchar(42),

    UNIQUE (order_id, src_chain),
    FOREIGN KEY (sell_token) REFERENCES tokens (id),
    FOREIGN KEY (buy_token) REFERENCES tokens (id)
);

CREATE TABLE match_orders
(
    id           bigserial PRIMARY KEY,
    match_id     bigint      NOT NULL,
    src_chain    bigint      NOT NULL,
    origin_order bigint      NOT NULL,
    order_id     bigint      NOT NULL,
    order_chain  bigint      NOT NULL,
    creator      varchar(42) NOT NULL,
    sell_token   bigint      NOT NULL,
    sell_amount  numeric(78) NOT NULL,
    state        smallint    NOT NULL,
    use_relayer  boolean     NOT NULL,

    UNIQUE (match_id, src_chain),
    FOREIGN KEY (sell_token) REFERENCES tokens (id),
    FOREIGN KEY (origin_order) REFERENCES orders (id),
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
DROP TABLE tokens;
DROP TABLE last_blocks;
