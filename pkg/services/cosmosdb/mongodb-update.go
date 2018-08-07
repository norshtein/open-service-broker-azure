package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	m *mongoAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"mongo account update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	m *mongoAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", m.updateReadRegions),
		service.NewUpdatingStep("waitForReadRegionsReady", m.waitForReadRegionsReady),
		service.NewUpdatingStep("updateARMTemplate", m.updateARMTemplate),
	)
}

func (m *mongoAccountManager) updateReadRegions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	up := instance.UpdatingParameters
	pp := instance.ProvisioningParameters

	goTemplateParameters, err := m.buildGoTemplateParamsOnlyRegionChanged(pp, up, dt, "MongoDB")
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	tags := getTags(up)

	_, err = m.armDeployer.Update(
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

func (m *mongoAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	up := instance.UpdatingParameters
	goTemplateParameters, err := m.buildGoTemplateParams(up, dt, "MongoDB")
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	tags := getTags(up)

	_, err = m.armDeployer.Update(
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
