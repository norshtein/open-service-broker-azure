package apimanagement

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2017-03-01/apimanagement"
)

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
		service.NewProvisioningStep("enableRESTAPI", s.enableRESTAPI),
	)
}

func (s *serviceManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error){
	dt := instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
	}
	dtMap,err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil{
		return nil, nil, err
	}

	pp := provisioningParameters{}
	if err := service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil{
		return nil, nil, err
	}

	_, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil,
		map[string]interface{}{
			"name": pp.apiName,
			"adminEmail" : pp.adminEmail,
			"orgName": pp.orgName,
			"tier": instance.Plan.GetProperties().Extended["tier"],
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	dt.apiName = pp.apiName
	dt.orgName = pp.orgName
	dt.adminEmail = pp.adminEmail

	dtMap,err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, nil
}

func (s *serviceManager) enableRESTAPI(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error){
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pp := provisioningParameters{}
	if err := service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil{
		return nil, nil, err
	}

	enabled := true
	_, err := s.tenantAccessClient.Update(ctx,
		instance.ResourceGroup,
		pp.apiName,
		apimanagement.AccessInformationUpdateParameters{
			Enabled: &enabled,
	},
		"*")
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
