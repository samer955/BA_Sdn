create table latency
(
    id      serial
        constraint latency_pk
            primary key,
    hostname text,
    node_id text,
    ip      text,
    latency integer,
    time    timestamp with time zone not null
);

alter table latency
    owner to "user";

create table cpu
(
    id      serial
        constraint cpu_pk
            primary key,
    hostname text,
    node_id text,
    ip      text,
    usage   integer,
    time    timestamp with time zone not null
);

alter table cpu
    owner to "user";

create table ram
(
    id      serial
        constraint ram_pk
            primary key,
    hostname text,
    node_id text,
    ip      text,
    usage   integer,
    time    timestamp with time zone not null
);

alter table ram
    owner to "user";
