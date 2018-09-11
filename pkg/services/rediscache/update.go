package rediscache

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters

	// Can't update the instance from a larger capacity to a smaller capacity
	provisionCapacity := pp.GetInt64("skuCapacity")
	updateCapacity := up.GetInt64("skuCapacity")
	if provisionCapacity > updateCapacity {
		return fmt.Errorf("can not update an instance from larger capacity %d to"+
			"smaller capacity %d", provisionCapacity, updateCapacity)
	}

	// Can't update `shardCount` and `skuCapacity` at the same time
	if strings.ToLower(instance.Plan.GetName()) == premium {
		if up.GetInt64("skuCapacity") != 0 && up.GetInt64("shardCount") != 0 {
			return fmt.Errorf("can not update `shardCount` and `skuCapacity` at the same time")
		}
	}
	return nil
}

func (s *serviceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *serviceManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	up := instance.UpdatingParameters
	tagsObj := up.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	_, err := s.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		buildGoTemplate(instance, update),
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating redis instance %s", err)
	}

	nonSSLPortEnabled := up.GetString("enableNonSslPort")
	if nonSSLPortEnabled != "" {
		if nonSSLPortEnabled == "enabled" {
			dt.NonSSLEnabled = true
		} else {
			dt.NonSSLEnabled = false
		}
	}
	return dt, nil
}
