package cosmosdb

import (
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	c *cosmosAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	c *cosmosAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

// This function will build a map in which only read regions changed.
func (c *cosmosAccountManager) buildGoTemplateParamsOnlyRegionChanged(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind
	p["location"] = pp.GetString("location")
	p["readRegions"] = pp.GetStringArray("readLocations")

	filters := []string{}
	ipFilters := pp.GetObject("ipFilters")
	if ipFilters.GetString("allowAzure") != disabled {
		filters = append(filters, "0.0.0.0")
	} else if ipFilters.GetString("allowPortal") != disabled {
		// Azure Portal IP Addresses per:
		// https://aka.ms/Vwxndo
		//|| Region            || IP address(es) ||
		//||=====================================||
		//|| China             || 139.217.8.252  ||
		//||===================||================||
		//|| Germany           || 51.4.229.218   ||
		//||===================||================||
		//|| US Gov            || 52.244.48.71   ||
		//||===================||================||
		//|| All other regions || 104.42.195.92  ||
		//||                   || 40.76.54.131   ||
		//||                   || 52.176.6.30    ||
		//||                   || 52.169.50.45   ||
		//||                   || 52.187.184.26  ||
		//=======================================||
		// Given that we don't really have context of the cloud
		// we are provisioning with right now, use all of the above
		// addresses.
		filters = append(filters,
			"104.42.195.92",
			"40.76.54.131",
			"52.176.6.30",
			"52.169.50.45",
			"52.187.184.26",
			"51.4.229.218",
			"139.217.8.252",
			"52.244.48.71",
		)
	} else {
		filters = append(filters, "0.0.0.0")
	}
	filters = append(filters, ipFilters.GetStringArray("allowedIPRanges")...)
	if len(filters) > 0 {
		p["ipFilters"] = strings.Join(filters, ",")
	}
	p["consistencyPolicy"] = pp.GetObject("consistencyPolicy").Data
	return p, nil
}
