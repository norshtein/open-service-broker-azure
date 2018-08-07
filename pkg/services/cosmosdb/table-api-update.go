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
		service.NewUpdatingStep("updateARMTemplate", t.updateARMTemplate),
		service.NewUpdatingStep("waitForReadRegionsReady", t.waitForReadRegionsReady),
	)
}

func (t *tableAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	up := instance.UpdatingParameters
	goTemplateParameters, err := t.buildGoTemplateParams(up, dt, "GlobalDocumentDB")
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	goTemplateParameters["capability"] = "EnableTable"
	tags := getTags(up)
	tags["defaultExperience"] = "Table"

	_, err = t.armDeployer.Update(
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
