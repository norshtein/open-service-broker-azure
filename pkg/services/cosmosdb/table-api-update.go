package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	t *tableAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"table account update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	t *tableAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", t.updateReadRegions),
	)
}
