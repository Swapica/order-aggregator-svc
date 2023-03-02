-- +migrate Up
ALTER TABLE orders
    ADD CONSTRAINT executed_by_match_fk FOREIGN KEY (executed_by_match) REFERENCES match_orders (id),
    ADD CONSTRAINT match_id_fk FOREIGN KEY (match_id, dest_chain) REFERENCES match_orders (match_id, src_chain);

-- +migrate Down
ALTER TABLE orders
    DROP CONSTRAINT executed_by_match_fk,
    DROP CONSTRAINT match_id_fk;
