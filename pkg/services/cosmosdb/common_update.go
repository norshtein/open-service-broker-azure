package cosmosdb

import (
	"fmt"

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

// This function is used in update. It will build a map in which
// only read regions changed. The rest will keep the same with provision parameter.
func (c *cosmosAccountManager) buildGoTemplateParamsOnlyRegionChanged(
	pp *service.ProvisioningParameters,
	up *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
) (map[string]interface{}, error) {
	readLocations := up.GetStringArray("readLocations")
	readLocations = append([]string{pp.GetString("location")}, readLocations...)
	return c.buildGoTemplateParamsCore(
		pp,
		dt,
		kind,
		readLocations,
	)
}
