package apimanagement

import (
	apiManagementSDK "github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2016-10-10/apimanagement" // nolint: 111
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer        arm.Deployer
	servicesClient     apiManagementSDK.ServicesClient
	tenantAccessClient apiManagementSDK.TenantAccessClient
}

func New(
	armDeployer arm.Deployer,
	servicesClient apiManagementSDK.ServicesClient,
	tenantAccessClient apiManagementSDK.TenantAccessClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:        armDeployer,
			servicesClient:     servicesClient,
			tenantAccessClient: tenantAccessClient,
		},
	}
}

func (m *module) GetName() string {
	return "apimanagement"
}
