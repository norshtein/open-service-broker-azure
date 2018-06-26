package apimanagement

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
}

type credentials struct {
	BaseURL     string `json:"baseUrl"`
	Identifier string `json:"identifier"`
	ExpiryDate string `json:"expiryDate"`
	Key     string `json:"key"`
}

func (
s *serviceManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"apiName",
			"adminEmail",
			"orgName",
		},
		//TODO: Add regular expression to the schema
		PropertySchemas: map[string]service.PropertySchema{
			"apiName": &service.StringPropertySchema{
				Description: "The name of the service, the api endpoint will be "+
					"<apiName>.azure-api.net",
			},
			"adminEmail": &service.StringPropertySchema{
				Description: "The administrator's email address to receive all " +
					"system notifications sent from API Management",
			},
			"orgName": &service.StringPropertySchema{
				Description: "Client ID (username) for an existing service principal," +
					"which will be granted access to the new vault.",
			},
		},
	}
}

func (d *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (d *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}