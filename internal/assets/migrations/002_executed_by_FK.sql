-- +migrate Up
ALTER TABLE orders
    ADD CONSTRAINT match_id_fk FOREIGN KEY (match_id, destination_chain) REFERENCES match_orders (match_id, src_chain);

-- +migrate Down
ALTER TABLE orders
    DROP CONSTRAINT match_id_fk;
