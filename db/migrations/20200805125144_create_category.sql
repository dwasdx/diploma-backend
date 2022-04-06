
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `sl_categories`
(
    `id`         INT          NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `sl_products`
(
    `id`          INT          NOT NULL AUTO_INCREMENT,
    `title`       VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
    `category_id` INT          NOT NULL,
    `created_at`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `sl_categories`;
DROP TABLE `sl_products`;
