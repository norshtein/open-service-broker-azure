# Disaster Recovery

This is a list of guidance for setting up disaster recovery for service instances in the service broker. Setting up disaster recovery can make your application more robust. Since unpredictable outage of data center and Azure region might happen, without disaster recovery, your application might be unavailable due to the disaster.

Currently, services supporting disaster recovery are listed below:

- [MSSQL](#mssql)

- [Cosmos DB](#cosmosdb)

 

## MSSQL

To be done by Zhongyi.



## Cosmos DB

Cosmos DB is a globally distributed database. To setting up disaster recovery for Cosmos DB, you should have a Cosmos DB instance which has at lease two regions. 

In the service broker, we implement disaster recovery for Cosmos DB by providing `-*registered` services. These services don't create any real resource in provisioning step. They are used to register existing Cosmos instances and provide credential of existing Cosmos instances. So when a disaster happens, you can use `*-registered` service to fetch the credentials in another environment.

Depending on your infrastructure, we have provided two kind of guidance for you:

- If you have (at least) two CF clusters, and you want to recover from an outage of Azure region and datacenter, see [this](cosmosdb-multi-region.md).
- If you have only one running environment ( CF, k8s, ...) , and you want to recover from an outage of datacenter, which causes your write region out of service, see [this](cosmosdb-one-region.md). **Note**: If there is an outage of Azure region and your whole running environment is unavailable, you can't recovery from the disaster and you have to wait until the Azure region comes online.