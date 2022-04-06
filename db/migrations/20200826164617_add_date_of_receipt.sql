
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `sl_item` ADD `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP AFTER `updated_at`;
ALTER TABLE `sl_item_list` ADD `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP AFTER `updated_at`;
ALTER TABLE `sl_shared_lists` ADD `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP AFTER `updated_at`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `sl_item` DROP `received_at`;
ALTER TABLE `sl_item_list` DROP `received_at`;
ALTER TABLE `sl_shared_lists` DROP `received_at`;
