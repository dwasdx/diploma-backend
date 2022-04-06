
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `sl_notifications` (
                                    `id` varchar(36) NOT NULL,
                                    `message` VARCHAR(512) NOT NULL,
                                    `user_id` VARCHAR(36) NOT NULL,
                                    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

ALTER TABLE `sl_notifications`
    ADD PRIMARY KEY (`id`),
    ADD KEY `user_id` (`user_id`,`created_at`);
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `sl_notifications`;
