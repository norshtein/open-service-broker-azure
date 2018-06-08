# [Azure API Management](https://azure.microsoft.com/en-us/services/api-management/)

## Services & Plans

### Service: azure-api-management

| Plan Name   | Description                 |
| ----------- | --------------------------- |
| `developer` | Developer Tier(No SLA)      |
| `basic`     | Basic Tier(99.9 SLA, %)     |
| `standard`  | Standard Tier(99.9 SLA, %)  |
| `premium`   | Premium Tier(99.95* SLA, %) |

#### Behaviors

##### Provision

Provisions a new API Management Service, and enable API Management REST API on it.

###### Provisioning Parameters

| Parameter Name  | Type                | Description                                                  | Required                                                     | Default Value                                                |
| --------------- | ------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `apiName`       | `string`            | The name of the api management service, should be unique on Azure. Your management url will be `<apiName> .management.azure-api.net` | Y                                                            |                                                              |
| `adminEmail`    | `string`            | The e-mail address to receive all system notifications sent from API Management | Y                                                            |                                                              |
| `orgName`       | `string`            | The name of your organization for use in the developer portal and e-mail notifications. | Y                                                            |                                                              |
| `location`      | `string`            | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured.                |
| `resourceGroup` | `string`            | The (new or existing) resource group with which to associate new resources. | N                                                            | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags`          | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N                                                            | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |



##### Bind

Generates and returns an access token, which expires after 1 month.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the folling connection details and credentials:

| Field Name   | Type     | Description                                                  |
| ------------ | -------- | ------------------------------------------------------------ |
| `baseUrl`    | `string` | Management URL, develops can send RESTful request to this URL. |
| `expiryDate` | `string` | The expriy date of the token.                                |
| `identifier` | `string` |                                                              |
| `key`        | `string` | Signature key used to create the access token.               |

Create your access token using the folling format:

``uid={identifier}&ex={expiryDate}&sn={key}``

Use these values to create an `Authorization` header in every request to the API Management REST API, as shown in the following example.

`Authorizaton: SharedAccessSignature uid=53dd860e1b72ff0467030003&ex=2014-08-04T22:03:00.0000000Z&sn=ItH6scUyCazNKHULKA0Yv6T+Skk4bdVmLqcPPPdWoxl2n1+rVbhKlplFrqjkoUFRr0og4wjeDz4yfThC82OjfQ== `



##### Unbind

Does nothing.



##### Deprovision

Deletes the API Management Service.