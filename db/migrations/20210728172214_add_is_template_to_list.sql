
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `sl_item_list` ADD `is_template` BOOLEAN NOT NULL DEFAULT FALSE AFTER `name`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `sl_item_list` DROP `is_template`;
