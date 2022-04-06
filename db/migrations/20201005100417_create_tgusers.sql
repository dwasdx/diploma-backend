-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `sl_tg_users`
(
    `tg_id`      int          NOT NULL,
    `username`   varchar(255) NOT NULL,
    `user_id`    varchar(36)  NOT NULL,
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;

ALTER TABLE `sl_tg_users`
    ADD PRIMARY KEY (`tg_id`);
COMMIT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `sl_tg_users`;
