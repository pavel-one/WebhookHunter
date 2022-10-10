alter table requests
    drop constraint fk_requests_hunters_id;

alter table requests
    drop column hunter_id;

create table channels
(
    id        bigserial
        primary key,
    hunter_id varchar
        constraint channels_hunters_id_fk
            references hunters
            on delete cascade,
    path      varchar not null,
    redirect  varchar
);

alter table requests
    add channel_id bigint;

alter table requests
    add constraint requests_channels_id_fk
        foreign key (channel_id) references channels (id)
            on delete cascade;