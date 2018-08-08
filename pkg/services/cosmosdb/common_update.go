package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) updateDeployment(
	pp *service.ProvisioningParameters,
	up *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
	capability string,
	additionalTags map[string]string,
) error {
	p, err := c.buildGoTemplateParams(up, dt, kind)
	if err != nil {
		return err
	}
	if capability != "" {
		p["capability"] = capability
	}
	tags := getTags(pp)
	for k, v := range additionalTags {
		tags[k] = v
	}
	err = c.deployUpdatedARMTemplate(
		up,
		dt,
		p,
		tags,
	)
	if err != nil {
		return fmt.Errorf("error deploying ARM template: %s", err)
	}
	return nil
}

func (c *cosmosAccountManager) updateReadRegions(
	pp *service.ProvisioningParameters,
	up *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
	capability string,
	additionalTags map[string]string,
) error {
	p, err := c.buildGoTemplateParamsOnlyRegionChanged(pp, up, dt, kind)
	if err != nil {
		return err
	}
	if capability != "" {
		p["capability"] = capability
	}
	tags := getTags(pp)
	for k, v := range additionalTags {
		tags[k] = v
	}
	err = c.deployUpdatedARMTemplate(
		up,
		dt,
		p,
		tags,
	)
	if err != nil {
		return fmt.Errorf("error deploying ARM template: %s", err)
	}
	return nil
}

func (c *cosmosAccountManager) deployUpdatedARMTemplate(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	goParams map[string]interface{},
	tags map[string]string,
) error {
	_, err := c.armDeployer.Update(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return fmt.Errorf("error deploying ARM template: %s", err)
	}
	return nil
}

// This function will build a map in which only read regions changed.
// This function is used in update.
func (c *cosmosAccountManager) buildGoTemplateParamsOnlyRegionChanged(
	provisioningParameters *service.ProvisioningParameters,
	updatingParameters *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind
	p["location"] = provisioningParameters.GetString("location")
	p["readRegions"] = updatingParameters.GetStringArray("readLocations")
	if provisioningParameters.GetString("autoFailoverEnabled") == enabled {
		p["enableAutomaticFailover"] = true
	} else {
		p["enableAutomaticFailover"] = false
	}

	filters := []string{}
	ipFilters := provisioningParameters.GetObject("ipFilters")
	if ipFilters.GetString("allowAzure") == disabled && ipFilters.GetString("allowPortal") != disabled {
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
	p["consistencyPolicy"] = provisioningParameters.GetObject("consistencyPolicy").Data
	return p, nil
}
