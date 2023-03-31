create table utils_script
(
    id          int identity
        constraint utils_script_pk
            primary key,
    name        varchar(255),
    script_name varchar(255) not null
);