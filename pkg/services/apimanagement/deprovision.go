package apimanagement

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", s.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteAPIManagementService",
			s.deleteAPIManagementService,
		),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters

	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (s *serviceManager) deleteAPIManagementService(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pp := instance.ProvisioningParameters
	if _, err := s.servicesClient.Delete(
		ctx,
		pp.GetString("resourceGroup"),
		pp.GetString("apiName"),
	); err != nil {
		return nil, fmt.Errorf("error deleting api management service: %s", err)
	}
	return instance.Details, nil
}
