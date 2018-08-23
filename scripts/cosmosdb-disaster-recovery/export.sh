#!/bin/bash

while getopts :o:h:a:d: option
do
	case "$option"
		in
		o) output_file_name=${OPTARG};;
		h) redis_host=${OPTARG};;
		a) redis_password=${OPTARG};;
		d) system_domain=${OPTARG};;
		\?) unknown_flag=${OPTARG};;
	esac
done

if [ ! -z $unknown_flag ]; then
	echo "Unknown flag $unknown_flag"
	exit 1
fi
if [ -z $redis_host ]; then
	echo 'Redis host is missed, use "-h <REDIS_HOST>" to specify the redis host'
	exit 1
fi
if [ -z $redis_password ]; then
	echo 'Redis password is missed, use "-a <REDIS_PASSWORD>" to specify the redis password'
	exit 1
fi
if [ -z $system_domain ]; then
	echo 'System domain is missed, use "-d <SYSTEM_DOMAIN>" to specify the system domain"'
	exit 1
fi
if [ -z $output_file_name ]; then
	echo 'Output file name is missed, use "-o <OUTPUT_FILE_NAME>" to specify the output file name'
	exit 1
fi

declare -A broker_plan_id_to_service_name=( ["58d7223d-934e-4fb5-a046-0c67781eb24e"]="azure-cosmosdb-sql" ["71168d1a-c704-49ff-8c79-214dd3d6f8eb"]="azure-cosmosdb-sql-account" ["c821c68c-c8e0-4176-8cf2-f0ca582a07a3"]="azure-cosmosdb-sql-database" ["86fdda05-78d7-4026-a443-1325928e7b02"]="azure-cosmosdb-mongo-account" ["126a2c47-11a3-49b1-833a-21b563de6c04"]="azure-cosmosdb-graph-account" ["c970b1e8-794f-4d7c-9458-d28423c08856"]="azure-cosmosdb-table-account" )
declare -A broker_plan_id_to_plan_name=( ["58d7223d-934e-4fb5-a046-0c67781eb24e"]="sql-api" ["71168d1a-c704-49ff-8c79-214dd3d6f8eb"]="account" ["c821c68c-c8e0-4176-8cf2-f0ca582a07a3"]="database" ["86fdda05-78d7-4026-a443-1325928e7b02"]="account" ["126a2c47-11a3-49b1-833a-21b563de6c04"]="account" ["c970b1e8-794f-4d7c-9458-d28423c08856"]="account" )
declare -A cf_plan_id_to_service_name
declare -A cf_plan_id_to_plan_name
declare -A service_instance_id_to_name
declare -A service_instance_id_to_cf_plan_id

echo "Getting token..."
TOKEN=$(cf oauth-token)

echo "Getting service plans..."
service_plans=$(curl -k "https://api.${system_domain}/v2/service_plans" -X GET -H "Authorization: ${TOKEN}")

echo "Processing plans..."
IFS=$'\n'
for row in $(echo $service_plans | jq '.resources' | jq -c '.[]'); do
	guid=$(echo $row | jq -r '.metadata.guid')
	unique_id=$(echo $row | jq -r '.entity.unique_id')
	for broker_id in "${!broker_plan_id_to_service_name[@]}"; do
		if [ $broker_id = $unique_id ]; then
			echo "Matched $guid: $unique_id"
			cf_plan_id_to_service_name[${guid}]=${broker_plan_id_to_service_name[${broker_id}]}
			cf_plan_id_to_plan_name[${guid}]=${broker_plan_id_to_plan_name[${broker_id}]}
			break
		fi
	done
done

echo "Getting service instances..."
service_instances=$(curl -k "https://api.${system_domain}/v2/service_instances" -X GET -H "Authorization: ${TOKEN}")

echo "Processing service instances..."
for row in $(echo $service_instances | jq '.resources' | jq -c '.[]'); do
	service_instance_guid=$(echo $row | jq -r '.metadata.guid')
	service_plan_guid=$(echo $row | jq -r '.entity.service_plan_guid')
	service_instance_name=$(echo $row | jq -r '.entity.name')
	for plan_id in "${!cf_plan_id_to_service_name[@]}"; do
		if [ $plan_id = $service_plan_guid ]; then
			echo "Matched service instance name ${service_instance_name} having id ${service_instance_guid} with plan id ${service_plan_guid} with service name ${cf_plan_id_to_service_name[${service_plan_guid}]} with plan name ${cf_plan_id_to_plan_name[${service_plan_guid}]}"
			service_instance_id_to_name[${service_instance_guid}]=${service_instance_name}
			service_instance_id_to_cf_plan_id[${service_instance_guid}]=${service_plan_guid}
			break
		fi
	done
done

echo "Getting service bindings..."
service_bindings=$(curl -k "https://api.${system_domain}/v2/service_bindings" -X GET -H "Authorization: ${TOKEN}")
echo "Processing service bindings..."
result_json='['
for row in $(echo $service_bindings | jq '.resources' | jq -c '.[]'); do
	service_instance_guid=$(echo $row | jq -r '.entity.service_instance_guid')
	loop_flag="true"
	for instance_id in "${!service_instance_id_to_name[@]}"; do
		if [ $loop_flag != "true" ]; then
			break
		fi

		if [ $service_instance_guid = $instance_id ]; then
			instance_name=${service_instance_id_to_name[${service_instance_guid}]}
			instance_plan_id=${service_instance_id_to_cf_plan_id[${service_instance_guid}]}
			instance_plan_name=${cf_plan_id_to_plan_name[${instance_plan_id}]}
			instance_service_name=${cf_plan_id_to_service_name[${instance_plan_id}]}

			accountname=""
			username=$(echo $row | jq -r '.entity.credentials.username // empty')
			if [ ! -z $username ]; then
				account_name=$username
			else
				uri=$(echo $row | jq -r '.entity.credentials.uri')
				account_name=$(echo $uri | sed -r 's/^https:\/\/(.*?)\.documents.*/\1/')
			fi

			database_name=""
			if [ ${cf_plan_id_to_plan_name[${instance_plan_id}]} != "account" ]; then
				database_name=$(echo $row | jq -r '.entity.credentials.databaseName')
			fi
			
			resource_group_name=$(redis-cli -h ${redis_host} -a ${redis_password} get "instances:${instance_id}" | jq -r '.provisioningParameters.resourceGroup')
			provisioning_read_region=$(redis-cli -h ${redis_host} -a ${redis_password} get "instances:${instance_id}" | jq -r '.provisioningParameters.readRegions[0]')
			updating_read_region=$(redis-cli -h ${redis_host} -a ${redis_password} get "instances:${instance_id}" | jq -r '.updatingParameters.readRegions[0]')	
			if [ ${provisioning_read_region} = "null" ] && [ ${updating_read_region} = "null" ]; then
				loop_flag="false"
				break
			fi
			echo "Export service instance ${instance_name}."
			
			# echo "sn: ${instance_service_name}, pn: ${instance_plan_name}, in: ${instance_name}, an: ${account_name}, dn: ${database_name}"
			current_json=$(jq -n --arg rn "$resource_group_name" --arg sn "$instance_service_name" --arg pn "$instance_plan_name" --arg in "$instance_name" --arg an "$account_name" --arg dn "$database_name" '{resourceGroup: $rn, serviceName: $sn, planName: $pn, instanceName: $in, accountName: $an, databaseName: $dn}')
			if [ "$result_json" != '[' ]; then
				result_json="${result_json},"		
			fi
			result_json="${result_json}${current_json}"
			loop_flag="false"
		fi
	done
done
result_json="${result_json}]"
echo $result_json > ${output_file_name}
echo "Done."
