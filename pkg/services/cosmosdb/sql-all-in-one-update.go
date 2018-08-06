package cosmosdb

import (
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
		service.NewUpdatingStep("updateReadRegions", s.updateReadRegions),
	)
}
