create table admins
(
    id       bigserial
        primary key,
    login    varchar unique,
    password varchar,
    created_at timestamp without time zone
);