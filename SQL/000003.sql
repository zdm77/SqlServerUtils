alter table utills_catalog_fields
    add link_table_id int;



alter table utills_catalog_fields
    add is_user_create bit;

alter table utills_catalog_fields
    add is_user_modify bit;

alter table utills_catalog_fields
    add is_date_create bit;


alter table utills_catalog_fields
    add is_date_modify bit;

alter table utils_catalog_list
    add type_entity varchar(20);

update utils_catalog_list set type_entity = 'Справочники' where 1=1;

alter table utills_catalog_fields
    add is_access_check bit;

update utills_catalog_fields set is_access_check=0 where 1=1;

alter table utills_catalog_fields
    add is_foreign_field bit;

update utills_catalog_fields set is_foreign_field=0 where 1=1;

alter table utills_catalog_fields
    add order_by int;

alter table utills_catalog_fields
    add order_by_form int;


alter table utils_catalog_list
    add order_by_default varchar(40);

alter table utils_catalog_list
    add order_by_default_asc varchar(5);

update utils_catalog_list set order_by_default='id', order_by_default_asc='asc' where type_entity='Справочники';
update utils_catalog_list set order_by_default='id', order_by_default_asc='desc' where type_entity='Задачи';