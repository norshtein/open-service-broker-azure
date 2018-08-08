package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	g *graphAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"graph account update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	g *graphAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		// The cosmosDB has a contraint: it cannot update properties and add/remove regions at the same time,
		// so we must deal with the update twice, one time updating region, one time updating properties.
		service.NewUpdatingStep("updateReadRegions", g.updateReadRegions),
		service.NewUpdatingStep("waitForReadRegionsReady", g.waitForReadRegionsReady),
		service.NewUpdatingStep("updateARMTemplate", g.updateARMTemplate),
	)
}

func (g *graphAccountManager) updateReadRegions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := g.cosmosAccountManager.updateReadRegions(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableGremlin",
		map[string]string{
			"defaultExperience": "Graph",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (g *graphAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := g.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableGremlin",
		map[string]string{
			"defaultExperience": "Graph",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
