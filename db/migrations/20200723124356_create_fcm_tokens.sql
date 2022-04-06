
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `sl_fcm_tokens` (
                              `token` varchar(255) NOT NULL,
                              `platform` varchar(10) NOT NULL,
                              `user_id` varchar(36) NOT NULL,
                              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `sl_fcm_tokens`
    ADD KEY `user_id` (`user_id`),
    ADD KEY `token` (`token`);
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `sl_fcm_tokens`;