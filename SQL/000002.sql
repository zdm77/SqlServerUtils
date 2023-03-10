alter table utils_task_catalog
    add sheet_number int DEFAULT 1
go

update utils_task_catalog set sheet_number = 1 where 1=1
go