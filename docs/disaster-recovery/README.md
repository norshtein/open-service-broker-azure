# Disaster Recovery

This is a list of guidance for setting up disaster recovery for Cloud Foundry service instances. This guidance is suitable for you if:

**Scenario 1**

- You have two Cloud Foundry cluster in different Azure region. 
- You want to recover from an outage of data center where your service instance locates
- You want to recover from an outage of Azure region where your working Cloud Foundry cluster locates

or:

**Scenario 2**

- You have one Cloud Foundry cluster
- You want to recover from an outage of data center where your service instance locates



If you have set up disaster recovery successfully, you can:

- For both scenario 1 and scenario 2, recover from an outage of data center where your service instance locates. In this case, if you haven't set up disaster recovery, you can't use the service instance until the data center become available, which may cause your application cannot work during this time. With disaster recovery enabled, you can follow the guidance to start up the backup service instance, make your application use the backup instance, and therefore your application can continue to work ignoring the outage of the data center.
- For scenario 1, recover from an outage of Azure region where your Cloud Foundry cluster locates. In this case, if you haven't set up disaster recovery, you whole Cloud Foundry cluster will be unavailable, all of your application cannot be reached until the Azure region become online. With disaster recovery enabled, you can follow the guidance to boot up backup Cloud Foundry cluster, start up the backup service instance, make your application use the backup service instance, and therefore you will have a working Cloud Foundry cluster and your application can continue to work ignoring the outage of the Azure region.



Currently, services supporting disaster recovery are listed below:

- [Cosmos Database](./cosmos.md)