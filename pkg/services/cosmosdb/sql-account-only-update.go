package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	s *sqlAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("enableReadRegions", s.enableReadRegions),
	)
}
