-- +migrate Up
ALTER TABLE orders
    ADD CONSTRAINT executed_by_fk FOREIGN KEY (executed_by, dest_chain) REFERENCES match_orders (match_id, src_chain);

-- +migrate Down
ALTER TABLE orders
    DROP CONSTRAINT executed_by_fk;
