create table if not exists peer
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
    online_user integer,
    role        text,
    network     text
);

create unique index if not exists peer_uuid_uindex
    on peer (uuid);

create table if not exists status
(
    uuid     text not null
        constraint status_pk
            primary key,
    source   text,
    source_ip   text,
    target_ip   text,
    target      text,
    is_alive boolean,
    rtt      integer,
    time     timestamp with time zone
);

create unique index if not exists status_uuid_uindex
    on status (uuid);

create table if not exists ram
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

create unique index if not exists ram_uuid_uindex
    on ram (uuid);

create table if not exists cpu
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

create unique index if not exists cpu_uuid_uindex
    on cpu (uuid);

create table if not exists tcp
(
    uuid       text not null
        constraint tcp_pk
            primary key,
    hostname   text,
    ip         text,
    queue_size integer,
    received   integer,
    sent       integer,
    time       timestamp with time zone
);

create unique index if not exists tcp_uuid_uindex
    on tcp (uuid);

create table if not exists bandwidth
(
    uuid      text not null
        constraint bandwidth_pk
            primary key,
    id        text,
    source    text,
    target    text,
    total_in  bigint,
    total_out bigint,
    rate_in   integer,
    rate_out  integer,
    time      timestamp with time zone
);

create unique index if not exists bandwidth_uuid_uindex
    on bandwidth (uuid);