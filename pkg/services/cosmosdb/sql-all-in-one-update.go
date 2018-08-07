package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAllInOneManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"sql all in one update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	s *sqlAllInOneManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
		service.NewUpdatingStep("waitForReadRegionsReady", s.waitForReadRegionsReady),
	)
}

func (s *sqlAllInOneManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	up := instance.UpdatingParameters
	goTemplateParameters, err := s.buildGoTemplateParams(up, dt, "GlobalDocumentDB")
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	tags := getTags(up)
	tags["defaultExperience"] = "DocumentDB"

	_, err = s.armDeployer.Update(
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
