package cosmosdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *commonRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("fillInCommonCredentials", c.fillInCommonCredentials),                 // nolint:lll
		service.NewProvisioningStep("fillInDifferentiatedCredentials", c.fillInDifferentiatedCredentials), // nolint: lll
	)
}

func (c *commonRegisteredManager) fillInCommonCredentials(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	databaseAccountName := instance.ProvisioningParameters.GetString("accountName") // nolint: lll
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup") // nolint: lll

	databaseAccount, err := c.databaseAccountsClient.Get(
		ctx,
		resourceGroupName,
		databaseAccountName,
	)
	if err != nil {
		return nil, err
	}
	fqdn := *(databaseAccount.DatabaseAccountProperties.DocumentEndpoint)

	keys, err := c.databaseAccountsClient.ListKeys(
		ctx,
		resourceGroupName,
		databaseAccountName,
	)
	if err != nil {
		return nil, err
	}
	pk := *(keys.PrimaryMasterKey)

	dt := &cosmosdbInstanceDetails{}
	dt.DatabaseAccountName = databaseAccountName
	dt.PrimaryKey = service.SecureString(pk)
	dt.FullyQualifiedDomainName = fqdn
	dt.ConnectionString = service.SecureString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		),
	)
	return dt, nil
}

func (c *commonRegisteredManager) fillInDifferentiatedCredentials(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	serviceName := instance.Service.GetName()
	dt := instance.Details.(*cosmosdbInstanceDetails)

	switch serviceName {
	case mongoAccountRegistered:
		dt.ConnectionString = service.SecureString(
			fmt.Sprintf(
				"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
				dt.DatabaseAccountName,
				dt.PrimaryKey,
				dt.FullyQualifiedDomainName,
			),
		)
		// Allow to remove the https:// and the port 443 on the FQDN
		// This will allow to adapt the FQDN for Azure Public / Azure Gov ...
		// Before :
		// https://6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com:443
		// After :
		// 6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com
		hostnameNoHTTPS := strings.Join(
			strings.Split(dt.FullyQualifiedDomainName, "https://"),
			"",
		)
		dt.FullyQualifiedDomainName = strings.Join(
			strings.Split(hostnameNoHTTPS, ":443/"),
			"",
		)
	case tableAccountRegistered:
		dt.ConnectionString = service.SecureString(
			fmt.Sprintf(
				"DefaultEndpointsProtocol=https;AccountName=%s;"+
					"AccountKey=%s;TableEndpoint=%s",
				dt.DatabaseAccountName,
				dt.FullyQualifiedDomainName,
				dt.PrimaryKey,
			),
		)
	case sqlAccountRegistered:
	case graphAccountRegistered:
		dt.ConnectionString = service.SecureString(
			fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
				dt.FullyQualifiedDomainName,
				dt.PrimaryKey,
			),
		)
	default:
		return nil, fmt.Errorf("given service name %s is not vaild", serviceName)
	}
	return dt, nil
}
