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
		service.NewProvisioningStep("fillInInstanceDetails", c.fillInInstanceDetails), // nolint:lll
	)
}

func (c *commonRegisteredManager) fillInInstanceDetails(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pk, err := getPrimaryKey(
		ctx,
		instance.ProvisioningParameters,
		c.databaseAccountsClient,
	)
	if err != nil {
		return nil, err
	}

	dt := &cosmosdbInstanceDetails{}
	dt.DatabaseAccountName = instance.ProvisioningParameters.GetString("accountName") // nolint:lll
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

func getPrimaryKey(
	ctx context.Context,
	pp *service.ProvisioningParameters,
	databaseAccountClient cosmosSDK.DatabaseAccountsClient,
) (service.SecureString, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	accountName := pp.GetString("accountName")

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
