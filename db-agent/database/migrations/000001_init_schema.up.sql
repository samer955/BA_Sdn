create table peer
(
    uuid        text not null
        constraint peer_pk
            primary key,
    hostname    text,
    ip          text,
    os          text,
    platform    text,
    version     text,
    latency     integer,
    time        timestamp with time zone,
    node_id     text,
    online_user integer
);

alter table peer
    owner to "user";

create unique index peer_uuid_uindex
    on peer (uuid);

create table status
(
    uuid     text not null
        constraint status_pk
            primary key,
    source   text,
    target   text,
    is_alive boolean,
    rtt      integer,
    time     timestamp with time zone
);

alter table status
    owner to "user";

create unique index status_uuid_uindex
    on status (uuid);

create table ram
(
    uuid     text not null
        constraint ram_pk
            primary key,
    hostname text,
    node_id  text,
    ip       text,
    usage    integer,
    time     timestamp with time zone
);

alter table ram
    owner to "user";

create unique index ram_uuid_uindex
    on ram (uuid);

create table cpu
(
    uuid     text not null
        constraint cpu_pk
            primary key,
    node_id  text,
    ip       text,
    hostname text,
    model    text,
    usage    integer,
    time     timestamp with time zone
);

alter table cpu
    owner to "user";

create unique index cpu_uuid_uindex
    on cpu (uuid);

create table tcp
(
    uuid       text not null
        constraint tcp_pk
            primary key,
    hostname   text,
    ip         text,
    queue_size integer,
    received   integer,
    sent       integer
);

alter table tcp
    owner to "user";

create unique index tcp_uuid_uindex
    on tcp (uuid);

