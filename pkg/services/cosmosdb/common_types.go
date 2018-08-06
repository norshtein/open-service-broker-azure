package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type cosmosdbInstanceDetails struct {
	ARMDeploymentName        string               `json:"armDeployment"`
	DatabaseAccountName      string               `json:"name"`
	FullyQualifiedDomainName string               `json:"fullyQualifiedDomainName"`
	IPFilters                string               `json:"ipFilters"`
	ConnectionString         service.SecureString `json:"connectionString"`
	PrimaryKey               service.SecureString `json:"primaryKey"`
}

// cosmosCredentials encapsulates CosmosDB-specific details for connecting via
// a variety of APIs. This excludes MongoDB.
type cosmosCredentials struct {
	URI                     string `json:"uri"`
	PrimaryConnectionString string `json:"primaryConnectionString"`
	PrimaryKey              string `json:"primaryKey"`
}

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (
	c *cosmosAccountManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (c *cosmosAccountManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

// Allowed Cosmos DB read regions
var allowedReadRegions = map[string]bool{
	"westus2":            true,
	"westus":             true,
	"southcentralus":     true,
	"centraluseuap":      true,
	"centralus":          true,
	"northcentralus":     true,
	"canadacentral":      true,
	"eastus":             true,
	"eastus2euap":        true,
	"eastus2":            true,
	"canadaeast":         true,
	"northeurope":        true,
	"ukwest":             true,
	"uksouth":            true,
	"francecentral":      true,
	"westeurope":         true,
	"westindia":          true,
	"centralindia":       true,
	"southindia":         true,
	"southeastasia":      true,
	"eastasia":           true,
	"koreacentral":       true,
	"koreasouth":         true,
	"japaneast":          true,
	"japanwest":          true,
	"australiasoutheast": true,
	"australiaeast":      true,
}
