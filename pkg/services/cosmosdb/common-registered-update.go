package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	c *commonRegisteredManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return nil
}

func (
	c *commonRegisteredManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}
