ALTER TABLE transactions
    ADD COLUMN payload_ciphertext TEXT AFTER category_id,
    ADD COLUMN payload_nonce VARBINARY(32) AFTER payload_ciphertext,
    ADD COLUMN payload_tag VARBINARY(32) AFTER payload_nonce,
    DROP COLUMN notes,
    DROP COLUMN amount;
