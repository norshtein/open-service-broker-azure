# [Azure Cosmos DB](https://azure.microsoft.com/en-us/services/cosmos-db/)

_Note: This module is EXPERIMENTAL and is not included in the General Availability release of Open Service Broker for Azure. It will be added in a future OSBA release._

## Services & Plans

### Service: azure-cosmosdb-sql-registered

| Plan Name | Description                                             |
| --------- | ------------------------------------------------------- |
| `sql-api` | Database Account and Database configured to use SQL API |

#### Behaviors

##### Provision

Finds the CosmosDB database account and CosmosDB database with given name in service broker's subscription, and get credentials of account and database. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |
| `databaseName` | `string` | The (new or existing) resource group with which to associate new resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                 | Type     | Description                                                  |
| -------------------------- | -------- | ------------------------------------------------------------ |
| `uri`                      | `string` | The fully-qualified address and port of the CosmosDB database account. |
| `primaryKey`               | `string` | A secret key used for connecting to the CosmosDB database.   |
| `primaryConnectionString`  | `string` | The full connection string, which includes the URI and primary key. |
| `databaseName`             | `string` | The generated database name.                                 |
| `documentdb_database_id`   | `string` | The database name provided in a legacy key for use with Azure libraries. |
| `documentdb_host_endpoint` | `string` | The fully-qualified address and port of the CosmosDB database account provided in a legacy key for use with Azure libraries. |
| `documentdb_master_key`    | `string` | A secret key used for connecting to the CosmosDB database provided in a legacy key for use with Azure libraries. |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.

##### 

### Service: azure-cosmosdb-sql-account-registered

| Plan Name | Description                                |
| --------- | ------------------------------------------ |
| `account` | Database Account configured to use SQL API |

#### Behaviors

##### Provision

Finds the CosmosDB database account with given name in service broker's subscription, and get credentials of the account. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                | Type     | Description                                                  |
| ------------------------- | -------- | ------------------------------------------------------------ |
| `uri`                     | `string` | The fully-qualified address and port of the CosmosDB database account. |
| `primaryKey`              | `string` | A secret key used for connecting to the CosmosDB database.   |
| `primaryConnectionString` | `string` | The full connection string, which includes the URI and primary key. |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.



### Service: azure-cosmosdb-sql-database-registered

| Plan Name  | Description                                                  |
| ---------- | ------------------------------------------------------------ |
| `database` | Database on existing CosmosDB database account configured to use SQL API |

#### Behaviors

##### Provision

Finds the CosmosDB database account and CosmosDB database with given name in service broker's subscription, and get credentials of account and database. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |
| `databaseName` | `string` | The (new or existing) resource group with which to associate new resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                 | Type     | Description                                                  |
| -------------------------- | -------- | ------------------------------------------------------------ |
| `uri`                      | `string` | The fully-qualified address and port of the CosmosDB database account. |
| `primaryKey`               | `string` | A secret key used for connecting to the CosmosDB database.   |
| `primaryConnectionString`  | `string` | The full connection string, which includes the URI and primary key. |
| `databaseName`             | `string` | The generated database name.                                 |
| `documentdb_database_id`   | `string` | The database name provided in a legacy key for use with Azure libraries. |
| `documentdb_host_endpoint` | `string` | The fully-qualified address and port of the CosmosDB database account provided in a legacy key for use with Azure libraries. |
| `documentdb_master_key`    | `string` | A secret key used for connecting to the CosmosDB database provided in a legacy key for use with Azure libraries. |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.

### Service: azure-cosmosdb-mongo-account-registered

| Plan Name | Description                           |
| --------- | ------------------------------------- |
| `account` | MongoDB on Azure provided by CosmosDB |

#### Behaviors

##### Provision

Finds the CosmosDB database account with given name in service broker's subscription, and get credentials of the account. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name         | Type     | Description                                                  |
| ------------------ | -------- | ------------------------------------------------------------ |
| `host`             | `string` | The fully-qualified address of the CosmosDB database account. |
| `port`             | `int`    | The port number to connect to on the CosmosDB database account. |
| `username`         | `string` | The name of the database user.                               |
| `password`         | `string` | The password for the database user.                          |
| `connectionstring` | `string` | The full connection string, which includes the host, port, username, and password. |
| `uri`              | `string` | URI encoded string that represents the connection information |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.

### Service: azure-cosmosdb-graph-account-registered

| Plan Name | Description                                            |
| --------- | ------------------------------------------------------ |
| `account` | Database Account configured to use Graph (Gremlin) API |

#### Behaviors

##### Provision

Finds the CosmosDB database account with given name in service broker's subscription, and get credentials of the account. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                | Type     | Description                                                  |
| ------------------------- | -------- | ------------------------------------------------------------ |
| `uri`                     | `string` | The fully-qualified address and port of the CosmosDB database account. |
| `primaryKey`              | `string` | A secret key used for connecting to the CosmosDB database account. |
| `primaryConnectionString` | `string` | The full connection string, which includes the URI and primary key. |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.

### Service: azure-cosmosdb-table-account-registered

| Plan Name | Description                                  |
| --------- | -------------------------------------------- |
| `account` | Database Account configured to use Table API |

#### Behaviors

##### Provision

Finds the CosmosDB database account with given name in service broker's subscription, and get credentials of the account. The provision doesn't create any resource.

###### Provisioning Parameters

| Parameter Name | Type     | Description                                                  | Required | Default Value |
| -------------- | -------- | ------------------------------------------------------------ | -------- | ------------- |
| `accountName`  | `string` | The Azure region in which to provision applicable resources. | Y        |               |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                | Type     | Description                                                  |
| ------------------------- | -------- | ------------------------------------------------------------ |
| `uri`                     | `string` | The fully-qualified address and port of the CosmosDB database account. |
| `primaryKey`              | `string` | A secret key used for connecting to the CosmosDB database account. |
| `primaryConnectionString` | `string` | The full connection string, which includes the URI and primary key. |

##### Unbind

Does nothing.

##### Deprovision

Does nothing.