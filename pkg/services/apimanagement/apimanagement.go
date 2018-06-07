package apimanagement

import(
	apiManagementSDK "github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2017-03-01/apimanagement"       // nolint: 111
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	serviceClient apiManagementSDK.ServiceClient
	tenantAccessClient apiManagementSDK.TenantAccessClient
}

func New(
	armDeployer arm.Deployer,
	serviceClient apiManagementSDK.ServiceClient,
	tenantAccessClient apiManagementSDK.TenantAccessClient,
) service.Module{
	return &module{
		serviceManager: &serviceManager{
			armDeployer: armDeployer,
			serviceClient: serviceClient,
			tenantAccessClient: tenantAccessClient,
		},
	}
}

func (m *module) GetName() string {
	return "apimanagement"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}