create table utils_task
(
    id       int identity
        constraint utils_task_pk
            primary key,
    name     varchar(255),
    table_db varchar(100),
	str_header int default 1
	
)
go


create table settings_utils_task
(
    task_id     int
        constraint settings_utils_task_task_id_fk
            references utils_task (id) on
delete cascade,
    field_excel varchar(255),
    field_db    varchar(255),
	col_number int
)
go

