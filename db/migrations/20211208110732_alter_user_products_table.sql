-- +goose Up
-- +goose StatementBegin
ALTER TABLE sl_user_products MODIFY category_id INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sl_user_products MODIFY category_id VARCHAR(36);
-- +goose StatementEnd