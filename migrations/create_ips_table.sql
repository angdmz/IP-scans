CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

create table public.ips
(
    ip_from      bigint,
    ip_to        bigint,
    proxy_type   varchar(3),
    country_code char(2),
    country_name varchar(64),
    region_name  varchar(128),
    city_name    varchar(128),
    isp          varchar(256),
    domain       varchar(128),
    usage_type   varchar(11),
    asn          text,
    "as"         varchar(256),
    id           uuid default uuid_generate_v4() not null
        constraint ips_pk
            primary key
);


alter table public.ips
    add ip_from bigint;

alter table public.ips
    add ip_to bigint;

alter table public.ips
    add proxy_type varchar(3);

alter table public.ips
    add country_code char(2);

alter table public.ips
    add country_name varchar(64);

alter table public.ips
    add region_name varchar(128);

alter table public.ips
    add city_name varchar(128);

alter table public.ips
    add isp varchar(256);

alter table public.ips
    add domain varchar(128);

alter table public.ips
    add usage_type varchar(11);

alter table public.ips
    add asn text;

alter table public.ips
    add "as" varchar(256);

alter table public.ips
    add id uuid default uuid_generate_v4() not null;



