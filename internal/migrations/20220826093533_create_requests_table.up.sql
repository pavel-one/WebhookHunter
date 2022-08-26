create table requests
(
    id         bigserial                   not null
        constraint pk_requests
            primary key,
    request    json                        not null,
    hunter_id  varchar                     not null
        constraint fk_requests_hunters_id
            references hunters (id)
            on delete cascade,
    created_at timestamp without time zone not null
);

