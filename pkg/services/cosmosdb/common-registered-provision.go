package cosmosdb

import (
	"context"
	"fmt"
	"regexp"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *commonRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("fillInPKAndFQDN", c.fillInPKAndFQDN),               // nolint:lll
		service.NewProvisioningStep("fillInConnectionString", c.fillInConnectionString), // nolint: lll
	)
}

func (c *commonRegisteredManager) fillInPKAndFQDN(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	databaseAccountName := instance.ProvisioningParameters.GetString("accountName") // nolint: lll
	pk, err := getPrimaryKey(
		ctx,
		databaseAccountName,
		c.databaseAccountsClient,
	)
	if err != nil {
		return nil, err
	}

	dt := &cosmosdbInstanceDetails{}
	dt.DatabaseAccountName = databaseAccountName
	dt.PrimaryKey = pk
	dt.FullyQualifiedDomainName = fmt.Sprintf(
		"https://%s.documents.azure.com:443/",
		instance.ProvisioningParameters.GetString("accountName"),
	)
	dt.ConnectionString = service.SecureString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		),
	)
	return dt, nil
}

func (c *commonRegisteredManager) fillInConnectionString(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	serviceName := instance.Service.GetName()
	dt := instance.Details.(*cosmosdbInstanceDetails)

	switch serviceName {
	case "azure-cosmosdb-mongo-account-registered":
		dt.ConnectionString = service.SecureString(
			fmt.Sprintf(
				"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
				dt.DatabaseAccountName,
				dt.PrimaryKey,
				dt.FullyQualifiedDomainName,
			),
		)
	case "azure-cosmosdb-table-account-registered":
		dt.ConnectionString = service.SecureString(
			fmt.Sprintf(
				"DefaultEndpointsProtocol=https;AccountName=%s;"+
					"AccountKey=%s;TableEndpoint=%s",
				dt.DatabaseAccountName,
				dt.FullyQualifiedDomainName,
				dt.PrimaryKey,
			),
		)
	case "azure-cosmosdb-sql-account-registered":
	case "azure-cosmosdb-graph-account-registered":
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
func getPrimaryKey(
	ctx context.Context,
	accountName string,
	databaseAccountClient cosmosSDK.DatabaseAccountsClient,
) (service.SecureString, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	allAccounts, err := databaseAccountClient.List(ctx)
	if err != nil {
		return "", err
	}
	databaseAccount, err := findDatabaseAccountWithGivenName(
		allAccounts,
		accountName,
	)
	if err != nil {
		return "", err
	}

	resourceGroupName, err := extractResourceGroupFromID(
		*(databaseAccount.ID),
	)
	if err != nil {
		return "", err
	}

	keys, err := databaseAccountClient.ListKeys(
		ctx,
		resourceGroupName,
		accountName,
	)
	if err != nil {
		return "", err
	}

	pk := *(keys.PrimaryMasterKey)
	return service.SecureString(pk), nil
}

// Find database account with given name in the list, if none is found, return nil
func findDatabaseAccountWithGivenName(
	databaseAccountList cosmosSDK.DatabaseAccountsListResult,
	name string,
) (cosmosSDK.DatabaseAccount, error) {
	databaseAccounts := *(databaseAccountList.Value)
	for i := range databaseAccounts {
		databaseAccount := databaseAccounts[i]
		if *(databaseAccount.Name) == name {
			return databaseAccount, nil
		}
	}
	return cosmosSDK.DatabaseAccount{},
		fmt.Errorf("Given database account is not found in current subscription")
}

func extractResourceGroupFromID(
	id string,
) (string, error) {
	re, err := regexp.Compile(".*/resourceGroups/(.*?)/.*")
	if err != nil {
		return "", fmt.Errorf("Error compiling regular expression: %s", err)
	}
	res := re.FindStringSubmatch(id)

	if len(res) != 2 {
		return "", fmt.Errorf("Given id is not vaild")
	}
	return res[1], nil
}
