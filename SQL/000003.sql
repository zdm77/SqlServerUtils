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
alter table utils_catalog_list
    add p_key_name varchar(30);
	
update utils_catalog_list set p_key_name='id' where p_key_name='' or p_key_name is null;
alter table utills_catalog_fields
    add link_field_view varchar(50);
	
	
	alter table utills_catalog_fields
    add default 0 for precision;

alter table utills_catalog_fields
    add default 0 for scale;

alter table utills_catalog_fields
    add default 1 for is_nullable;


alter table utills_catalog_fields
    add default 0 for is_identity;


alter table utills_catalog_fields
    add default 0 for is_primary_key;


alter table utills_catalog_fields
    add default 1 for is_nullable_db;
	
	
	alter table utills_catalog_fields
    add default 1 for is_list;


alter table utills_catalog_fields
    add default 1 for is_form;


alter table utills_catalog_fields
    add default 0 for is_user_create;


alter table utills_catalog_fields
    add default 0 for is_user_modify;


alter table utills_catalog_fields
    add default 0 for is_date_create;


alter table utills_catalog_fields
    add default 0 for is_date_modify;


alter table utills_catalog_fields
    add default 0 for is_access_check;


alter table utills_catalog_fields
    add default 0 for is_foreign_field;


alter table utills_catalog_fields
    add default 0 for order_by;


alter table utills_catalog_fields
    add default 0 for order_by_form
go

alter table utills_catalog_fields
    add default '' for link_field_view
go

