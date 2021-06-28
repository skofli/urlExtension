create extension if not exists "uuid-ossp";
create table urls
(
    id           serial not null
        constraint urls_pk
            primary key,
    redirect_url text,
    long_url     text
);

alter table urls
    owner to postgres;

create unique index urls_id_uindex
    on urls (id);

