#!/usr/env/bin bash

count=0
while IFS=, read -r  city city_ascii lat lng country iso2 iso3 admin_name capital population id
do
	item="($country,$city_ascii,$lat,$lng)"
	((count++))
	echo $count,$item
	mysql -uroot <<SQLEOF
	use saga;
	insert into tbl_country_city (Country,City,Lat,Lng) values $item;
SQLEOF
	[[ $? -ne 0 ]] && { echo "write failed " && exit; }
done < $1
