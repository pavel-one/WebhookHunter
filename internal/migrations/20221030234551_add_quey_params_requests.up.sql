alter table requests
    add column path varchar not null default '',
    add column query json