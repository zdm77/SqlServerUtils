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

alter table utills_catalog_fields
    add is_access_check bit
go

update utills_catalog_fields set is_access_check=0 where 1=1;

alter table utills_catalog_fields
    add is_foreign_field bit
go

update utills_catalog_fields set is_foreign_field=0 where 1=1;

alter table utills_catalog_fields
    add order_by int
go

alter table utills_catalog_fields
    add order_by_form int
go
