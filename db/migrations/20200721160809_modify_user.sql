
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `sl_users` CHANGE `code` `code` VARCHAR(8) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `sl_users` CHANGE `code` `code` VARCHAR(8) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '';
