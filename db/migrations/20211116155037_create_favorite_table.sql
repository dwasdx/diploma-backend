-- +goose Up
-- +goose StatementBegin
create table sl_favorite
(
    id          varchar(36)                         NOT NULL,
    name        varchar(255)                        NOT NULL,
    owner_id    varchar(36) COLLATE utf8_general_ci NOT null,
    created_at  TIMESTAMP                           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    received_at TIMESTAMP                           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted  tinyint(1) NOT NULL DEFAULT '0',
    constraint sl_favorite_pk
        primary key (id),
    constraint sl_favorite_sl_users_id_fk
        foreign key (owner_id) references sl_users (id)
            on update cascade on delete cascade
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sl_favorite;
-- +goose StatementEnd
