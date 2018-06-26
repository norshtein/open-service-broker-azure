package apimanagement

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2016-10-10/apimanagement"
	"strings"
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
) (service.InstanceDetails, error){
	dt := &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
	}
	return dt, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters
	armTemplateParamters := map[string]interface{} {
		"name": pp.GetString("apiName"),
		"adminEmail": pp.GetString("adminEmail"),
		"orgName": pp.GetString("orgName"),
		"tier": instance.Plan.GetProperties().Extended["tier"],
	}
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	_, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		nil,
		armTemplateParamters,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (s *serviceManager) enableRESTAPI(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error){
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pp := instance.ProvisioningParameters
	enabled := true
	_, err := s.tenantAccessClient.Update(ctx,
		pp.GetString("resourceGroup"),
		pp.GetString("apiName"),
		apimanagement.AccessInformationUpdateParameters{
			Enabled: &enabled,
	},
		"*")

	// OSBA only provides an old version api-management go sdk, which treats http
	// response code 204 as an error, but in fact 204 indicates creating success.
	// So we add a special judgement here.
	if err != nil && !strings.Contains(err.Error(), "204"){
		return nil, err
	}
	return instance.Details, nil
}
