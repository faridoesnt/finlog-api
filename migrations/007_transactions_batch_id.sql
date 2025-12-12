ALTER TABLE transactions
    ADD COLUMN batch_id BIGINT NULL AFTER user_id,
    ADD INDEX idx_transactions_batch (batch_id),
    ADD CONSTRAINT fk_transactions_import_batch FOREIGN KEY (batch_id) REFERENCES import_batches(id)
        ON DELETE RESTRICT;
