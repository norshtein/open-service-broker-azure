package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	t *tableAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"table account update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	t *tableAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", t.updateReadRegions),
		service.NewUpdatingStep("waitForReadRegionsReady", t.waitForReadRegionsReady),
		service.NewUpdatingStep("updateARMTemplate", t.updateARMTemplate),
	)
}

func (t *tableAccountManager) updateReadRegions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := t.cosmosAccountManager.updateReadRegions(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableTable",
		map[string]string{
			"defaultExperience": "Table",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (t *tableAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := t.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableTable",
		map[string]string{
			"defaultExperience": "Table",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
