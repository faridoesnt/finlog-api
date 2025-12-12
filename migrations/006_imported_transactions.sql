CREATE TABLE IF NOT EXISTS imported_transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    batch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    payload_ciphertext TEXT NOT NULL,
    payload_nonce VARBINARY(32) NOT NULL,
    payload_tag VARBINARY(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_imported_transactions_batch FOREIGN KEY (batch_id) REFERENCES import_batches(id) ON DELETE CASCADE,
    CONSTRAINT fk_imported_transactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    KEY idx_imported_transactions_user (user_id),
    KEY idx_imported_transactions_batch (batch_id)
) ENGINE=InnoDB;
