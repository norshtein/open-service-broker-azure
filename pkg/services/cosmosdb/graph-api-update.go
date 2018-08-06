package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	g *graphAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadRegions(
		"graph account update",
		instance.UpdatingParameters.GetStringArray("readLocations"),
	)
}

func (
	g *graphAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadRegions", g.updateReadRegions),
	)
}
