
			select 
			tbl._id
			,tbl._date
			,tbl._shift
			,tbl._loader
			,tbl._kig_fact
			,tbl._kig_pass
			,tbl._rides
			,tbl._rides_ok
			,tbl._rides_low
			,tbl._rides_high
			,tbl._rides_fail
			,tbl.enterprise_id										_comp

			,c.work_name											_cargo
			,e.dir_for_plan											_comp_name
			,concat(i.[area], '-', i.[block])						_field
		
							from (select 
									trips_tmp_id											_id
									,cast(trip_date as date)								_date
									, enterprise_id
									, shovel_export_code
									,cargo_name
									,shift_n												_shift
									, trip_date



									,[shovel_export_code]									_loader	
									,[weight] / nullif([weight_mdl],0)*100				_kig_fact
									,[weight] / nullif([weight_std],0)*100				_kig_pass
									,1														_rides
									,case 
										when [weight] / nullif([weight_mdl],0) >= 0.95
											and [weight] / nullif([weight_mdl],0) <= 1.02
										then 1
										else 0 end											_rides_ok
									,case 
										when [weight] / nullif([weight_mdl],0) < 0.95
										then 1
										else 0 end											_rides_low
									,case 
										when [weight] / nullif([weight_mdl],0) > 1.02
										then 1
										else 0 end											_rides_high
									,case 
										when [weight] / nullif([weight_mdl],0) < 0.95
											and [weight] / nullif([weight_std],0) >= 0.95
											and [weight] / nullif([weight_std],0) <= 1.02
										then 1 else 0 end									_rides_fail
								from [PAC].[dbo].[trips_tmp] 
									-- фильтрация строк для отчета КИГ. Согласовано Аршинцев и Кавтарашвили
										where kig between 50 and 130
											and weight_mdl>=130
											and weight_std>0
											and weight_std is not null
											and distance_loaded between 0.1 and 30
											and speed_loaded between 5 and 60
											and speed_unloaded between 5 and 60) tbl

								left join [PAC].[dbo].[enterprise] e
									on tbl.enterprise_id = e.enterprise_id

								left join [dbo].[interval_pf] i 
									on tbl.shovel_export_code = i.vehicle_export_code
									and tbl.trip_date >= i.start_dat
									and tbl.trip_date <= i.end_dat

								left join [dbo].[vehicle_tmp] v 
									on v.vehicle_export_code = tbl.shovel_export_code

								left join [dbo].[cargo_type_tmp] c 
									on c.cargo_name = tbl.cargo_name -- enterprise id удалили, будет один справочник на всех