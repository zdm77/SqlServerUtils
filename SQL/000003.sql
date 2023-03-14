alter table utills_catalog_fields
    add link_table_id int
go


alter table utills_catalog_fields
    add is_user_create bit
go

alter table utills_catalog_fields
    add is_user_modify bit
go

alter table utills_catalog_fields
    add is_date_create bit
go



alter table utills_catalog_fields
    add is_date_modify bit
go

alter table utils_catalog_list
    add type_entity varchar(20)
go

update utils_catalog_list set type_entity = 'Справочники' where 1=1;