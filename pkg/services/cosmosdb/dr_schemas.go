package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func generateCommonRegisteredProvisioningParamsSchema() service.InputParametersSchema {
	propertySchemas := map[string]service.PropertySchema{
		"accountName": &service.StringPropertySchema{
			Title:       "Account Name",
			Description: "The database account name of existing instance",
		},
	}
	return service.InputParametersSchema{
		RequiredProperties: []string{"accountName"},
		PropertySchemas:    propertySchemas,
	}
}

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
	}
	return service.InputParametersSchema{
		RequiredProperties: []string{"accountName"},
		PropertySchemas:    propertySchemas,
	}
}
