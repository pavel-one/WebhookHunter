create table tokens
(
    id       bigserial
        primary key,
    token    varchar unique,
    created_at timestamp without time zone
);