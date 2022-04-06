-- +goose Up
-- +goose StatementBegin
ALTER TABLE sl_users MODIFY phone VARCHAR(15);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sl_users MODIFY phone INT;
-- +goose StatementEnd
