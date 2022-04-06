-- +goose Up
-- +goose StatementBegin
ALTER TABLE sl_user_products
    ADD `category_id`       varchar(36) COLLATE utf8_general_ci NOT NULL AFTER `owner_id`,
    ADD `global_product_id` int                                 NULL AFTER `category_id`,
    ADD `is_favorite`       tinyint(1)                          NOT NULL DEFAULT '0' AFTER `is_deleted`;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sl_user_products
    DROP category_id,
    DROP global_product_id,
    DROP is_favorite;
-- +goose StatementEnd