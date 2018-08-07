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
		service.NewUpdatingStep("updateARMTemplate", g.updateARMTemplate),
		service.NewUpdatingStep("waitForReadRegionsReady", g.waitForReadRegionsReady),
	)
}

func (g *graphAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	up := instance.UpdatingParameters
	goTemplateParameters, err := g.buildGoTemplateParams(up, dt, "GlobalDocumentDB")
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	goTemplateParameters["capability"] = "EnableGremlin"
	tags := getTags(up)
	tags["defaultExperience"] = "Graph"

	_, err = g.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.UpdatingParameters.GetString("resourceGroup"),
		instance.UpdatingParameters.GetString("location"),
		armTemplateBytes,
		goTemplateParameters,
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, err
}
