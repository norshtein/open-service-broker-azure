package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// nolint: lll
func generateCommonRegisteredProvisioningParamsSchema() service.InputParametersSchema {
	propertySchemas := map[string]service.PropertySchema{
		"accountName": &service.StringPropertySchema{
			Title:       "Account Name",
			Description: "The database account name of existing instance",
		},
		"resourceGroup": &service.StringPropertySchema{
			Title:       "Resource group",
			Description: "The resource group of the database account.",
		},
	}
	return service.InputParametersSchema{
		RequiredProperties: []string{"accountName", "resourceGroup"},
		PropertySchemas:    propertySchemas,
	}
}

// nolint: lll
func generateAllInOneRegisteredProvisioningParamsSchema() service.InputParametersSchema {
	propertySchemas := map[string]service.PropertySchema{
		"accountName": &service.StringPropertySchema{
			Title:       "Account Name",
			Description: "The database account name of existing instance",
		},
		"databaseName": &service.StringPropertySchema{
			Title:       "Database Name",
			Description: "The database name of existing instance",
		},
		"resourceGroup": &service.StringPropertySchema{
			Title:       "Resource group",
			Description: "The resource group of the database account.",
		},
	}
	return service.InputParametersSchema{
		RequiredProperties: []string{"accountName", "databaseName", "resourceGroup"},
		PropertySchemas:    propertySchemas,
	}
}
