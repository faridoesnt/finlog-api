ALTER TABLE categories
    ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1 AFTER icon_key;
