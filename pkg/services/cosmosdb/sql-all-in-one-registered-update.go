package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	s *sqlAllInOneRegisteredManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return nil
}

func (
	s *sqlAllInOneRegisteredManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}
