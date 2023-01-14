create table public.ips
(
    ip_from      integer,
    ip_to        integer,
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
    owner to admin;

