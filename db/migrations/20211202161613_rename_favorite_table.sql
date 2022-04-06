-- +goose Up
-- +goose StatementBegin
RENAME TABLE `sl_favorite` TO `sl_user_products`;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `sl_user_products` RENAME `sl_favorite`;
-- +goose StatementEnd
