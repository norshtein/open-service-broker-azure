package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	g *graphAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	g *graphAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", g.updateReadRegions),
	)
}
