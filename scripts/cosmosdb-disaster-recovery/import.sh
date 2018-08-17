#!/bin/bash

import_file_name=$1
for row in $(cat ${import_file_name} | jq -c '.[]' ); do
	service_name=$(echo $row | jq -r '.serviceName')
	plan_name=$(echo $row | jq -r '.planName')
	instance_name=$(echo $row | jq -r '.instanceName')
	account_name=$(echo $row | jq -r '.accountName')
	database_name=$(echo $row | jq -r '.databaseName')

	service_name="${service_name}-registered"
	instance_name="${instance_name}-registered"
	# echo $service_name  $plan_name $instance_name $account_name $database_name
	if [ $plan_name != "account" ]; then
		cf create-service "${service_name}" "${plan_name}" "${instance_name}" -c "{\"accountName\": \"${account_name}\", \"databaseName\": \"${database_name}\"}"
	else
		cf create-service "${service_name}" "${plan_name}" "${instance_name}" -c "{\"accountName\": \"${account_name}\"}"
	fi

done
