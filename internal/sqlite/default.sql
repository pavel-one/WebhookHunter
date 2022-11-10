PRAGMA foreign_keys = ON;

create table channels
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    path       varchar                           NOT NULL,
    created_at timestamp
);

create table requests
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    request    json                              NOT NULL,
    created_at timestamp                         NOT NULL,
    channel_id INTEGER                           NOT NULL,
    headers    json                              NOT NULL,
    path       varchar                           NOT NULL,
    query      json                              NOT NULL,
    FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);
