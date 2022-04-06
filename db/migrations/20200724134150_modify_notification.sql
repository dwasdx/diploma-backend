
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `sl_notifications`
    ADD `type` SMALLINT NULL DEFAULT NULL AFTER `id`,
    ADD `user_phone` VARCHAR(50) NULL DEFAULT NULL AFTER `user_id`,
    ADD `list_id` VARCHAR(36) NULL DEFAULT NULL AFTER `user_phone`,
    ADD `item_id` VARCHAR(36) NULL DEFAULT NULL AFTER `list_id`,
    ADD `target_user_id` VARCHAR(36) NOT NULL AFTER `item_id`;

ALTER TABLE `sl_notifications` CHANGE `user_id` `user_id` VARCHAR(36) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL;

ALTER TABLE `sl_notifications` DROP INDEX `user_id`;

ALTER TABLE `sl_notifications`
    ADD KEY `target_user_id` (`target_user_id`,`created_at`);
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `sl_notifications`
    DROP `type`,
    DROP `user_phone`,
    DROP `list_id`,
    DROP `item_id`,
    DROP `target_user_id`;

ALTER TABLE `sl_notifications` DROP INDEX `target_user_id`;

ALTER TABLE `sl_notifications`
    ADD KEY `user_id` (`user_id`,`created_at`);
COMMIT;