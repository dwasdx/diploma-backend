-- +goose Up
-- +goose StatementBegin
ALTER TABLE sl_favorite MODIFY COLUMN name VARCHAR(255)
    CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sl_favorite MODIFY COLUMN name VARCHAR(255)
    CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL
-- +goose StatementEnd
