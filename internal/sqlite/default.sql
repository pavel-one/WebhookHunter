create table channels
(
    id          bigserial
        primary key,
    hunter_slug varchar,
    path        varchar not null,
    redirect    varchar,
    created_at  timestamp
);

create table requests
(
    id         bigserial
        constraint pk_requests
            primary key,
    request    json      not null,
    created_at timestamp not null,
    channel_id bigint
        constraint requests_channels_id_fk
            references channels
            on delete cascade,
    headers    json,
    path       varchar   not null,
    query      json
);
