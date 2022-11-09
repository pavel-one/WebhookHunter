create table channels
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    path        varchar                           NOT NULL,
    created_at  timestamp
);

create table requests
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    request    json                              NOT NULL,
    created_at timestamp                         NOT NULL,
    channel_id INTEGER
        constraint requests_channels_id_fk
            references channels
            on delete cascade,
    headers    json                              NOT NULL,
    path       varchar                           NOT NULL,
    query      json                              NOT NULL
);
