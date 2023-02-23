create table utils_task_catalog
(
    id       int identity
        constraint utils_task_pk
            primary key,
    name     varchar(255),
    table_db varchar(100),
	str_header int default 1
	
)
go


create table utils_settings_task
(
    task_id     int
        constraint settings_utils_task_task_id_fk
            references utils_task_catalog (id) on
delete cascade,
    field_excel varchar(255),
    field_db    varchar(255),
	col_number int,
	field_type varchar(10)
)
go

create table utils_task
(
    id        int identity
        constraint utils_task_pk
            primary key,
    task_id   int          not null
        constraint utils_task_utils_task_catalog_id_fk
            references utils_task_catalog,
    name      varchar(100) not null,
    date_task date default GETDATE(),
    time_task time default GETDATE(),
    user_id   varchar(40) default user_name()
)
go


create table utils_catalog_list
(
    id         int identity
        constraint utils_catalog_list_pk
            primary key,
    name       varchar(255) not null,
    table_name varchar(255) not null
)
go




