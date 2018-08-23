package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *commonRegisteredManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("unregisterDatabaseAccount", c.unregisterDatabaseAccount), // nolint: lll
	)
}

func (c *commonRegisteredManager) unregisterDatabaseAccount(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	// do nothing, just for the framework to get the first step as it is required
	return instance.Details, nil
}
