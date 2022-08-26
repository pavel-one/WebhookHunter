create table hunters
(
    id         varchar                     not null
        constraint pk_hunters
            primary key,
    ip         varchar                     not null,
    created_at timestamp without time zone not null
);
