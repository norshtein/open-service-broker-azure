package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	t *tableAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	t *tableAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", t.updateReadRegions),
	)
}
