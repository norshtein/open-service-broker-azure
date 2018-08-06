package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return readRegionsValidator(
		"sql account update",
		[]interface{}{instance.UpdatingParameters.GetStringArray("readLocations")},
	)
}

func (
	s *sqlAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", s.updateReadRegions),
	)
}
