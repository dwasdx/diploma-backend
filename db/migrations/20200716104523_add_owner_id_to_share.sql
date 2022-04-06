
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `sl_shared_lists` ADD `owner_id` VARCHAR(36) NOT NULL AFTER `to_user_id`;
ALTER TABLE `sl_shared_lists` DROP PRIMARY KEY;
ALTER TABLE `sl_shared_lists` ADD PRIMARY KEY (`id`, `owner_id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `sl_shared_lists` DROP PRIMARY KEY;
ALTER TABLE `sl_shared_lists` ADD PRIMARY KEY (`id`);
ALTER TABLE `sl_shared_lists` DROP `owner_id`;
