create table info
(
    id       serial
        constraint latency_pk
            primary key,
    hostname text,
    node_id  text,
    ip       text,
    latency  integer,
    time     timestamp with time zone not null
);

alter table info
    owner to "user";

grant select on info to grafanareader;

create table cpu
(
    id       serial
        constraint cpu_pk
            primary key,
    hostname text,
    node_id  text,
    ip       text,
    usage    integer,
    time     timestamp with time zone not null
);

alter table cpu
    owner to "user";

grant select on cpu to grafanareader;

create table ram
(
    id       serial
        constraint ram_pk
            primary key,
    hostname text,
    node_id  text,
    ip       text,
    usage    integer,
    time     timestamp with time zone not null
);

alter table ram
    owner to "user";

grant select on ram to grafanareader;

create table process
(
    id       serial
        constraint process_pk
            primary key,
    name     text,
    cpu      double precision,
    ip       text,
    hostname text,
    time     timestamp with time zone
);

alter table process
    owner to "user";

create table status
(
    id        serial
        constraint status_pk
            primary key,
    source_id text,
    target_id text,
    is_alive  boolean,
    rtt       integer,
    time      timestamp with time zone
);

alter table status
    owner to "user";