# Disaster Recovery

This is a list of guidance for setting up disaster recovery for service instances in the service broker. Setting up disaster recovery can make your application more robust. Since unpredictable outage of data center and Azure region might happen, without disaster recovery, you may suffer:

- When there is an outage in data center where your service instance locates, you can't use the service instance until the data center recover from the disaster, which may cause your application cannot work during this time. 
- When there is an outage in Azure region where your Cloud Foundry clusters locates, you whole Cloud Foundry cluster will be unavailable, and all of your application cannot be reached until the Azure region become online. 



This guidance is suitable for you if:

**Scenario 1**

- You have two Cloud Foundry cluster in different Azure regions
- You want to recover from an outage of data center where your service instance locates
- You want to recover from an outage of Azure region where your working Cloud Foundry cluster locates

or:

**Scenario 2**

- You have one Cloud Foundry cluster
- You want to recover from an outage of data center where your service instance locates



If you have set up disaster recovery successfully, you can:

- For both scenario 1 and scenario 2, recover from an outage of data center where your service instance locates.  You can follow the guidance to start up the backup service instance, bind the back up service instance to your application, and therefore your application can continue to work ignoring the outage of the data center.
- For scenario 1, recover from an outage of Azure region where your Cloud Foundry cluster locates. With disaster recovery enabled, you can follow the guidance to boot up backup Cloud Foundry cluster, start up the backup service instance, bind the backup service instance to your application, and therefore you will have a working Cloud Foundry cluster and your application can continue to work ignoring the outage of the Azure region.



Currently, services supporting disaster recovery are listed below:

- [CosmosDB](./cosmosdb.md)