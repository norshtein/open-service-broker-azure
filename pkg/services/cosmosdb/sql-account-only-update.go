package cosmosdb

import (
	"context"
	"fmt"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	s *sqlAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("enableReadRegions", s.updateReadRegions),
	)
}

func (s *sqlAccountManager) updateReadRegions(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dt := instance.Details.(*cosmosdbInstanceDetails)
	var readLocations []string
	// If parameter `readLocations` is not specified, directlly return
	if readLocations = instance.UpdatingParameters.GetStringArray("readLocations"); len(readLocations) == 0 {
		return dt, nil
	}

	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := s.databaseAccountsClient
	nowDatabaseAccount, err := databaseAccountClient.Get(
		ctx,
		resourceGroupName,
		dt.DatabaseAccountName,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching created account's information: %s", err)
	}

	// Build new property
	nowProperties := nowDatabaseAccount.DatabaseAccountProperties
	newProperties := cosmosSDK.DatabaseAccountCreateUpdateProperties{}
	newProperties.ConsistencyPolicy = nowProperties.ConsistencyPolicy
	daot := string(nowProperties.DatabaseAccountOfferType)
	newProperties.DatabaseAccountOfferType = &daot
	newProperties.IPRangeFilter = nowProperties.IPRangeFilter
	newProperties.EnableAutomaticFailover = nowProperties.EnableAutomaticFailover
	newProperties.Capabilities = nowProperties.Capabilities

	// Wrap user provided read region
	locations := []cosmosSDK.Location{}
	locations = append(locations, contructLocation(accountName, instance.ProvisioningParameters.GetString("location"), 0))
	for i := range readLocations {
		locations = append(locations, contructLocation(accountName, readLocations[i], int32(i+1)))
	}
	newProperties.Locations = &locations

	// Build new parameter
	newParameter := cosmosSDK.DatabaseAccountCreateUpdateParameters{}
	newParameter.Kind = nowDatabaseAccount.Kind
	newParameter.ID = nowDatabaseAccount.ID
	newParameter.Name = nowDatabaseAccount.Name
	newParameter.Type = nowDatabaseAccount.Type
	newParameter.Location = nowDatabaseAccount.Location
	newParameter.Tags = nowDatabaseAccount.Tags
	newParameter.DatabaseAccountCreateUpdateProperties = &newProperties

	fmt.Println("Begin creating read regions")
	_, err = databaseAccountClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		dt.DatabaseAccountName,
		newParameter,
	)
	if err != nil {
		return nil, err
	}
	err = waitForRegionCreationCompletion(ctx, databaseAccountClient, resourceGroupName, accountName)
	return dt, nil
}
