package apimanagement

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ApiName string `json:"apiName"`
	AdminEmail string `json:"adminEmail"`
	OrgName string `json:"orgName"`
}

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ApiName string `json:"apiName"`
	AdminEmail string `json:"adminEmail"`
	OrgName string `json:"orgName"`
}

type credentials struct {
	BaseURL     string `json:"baseUrl"`
	Identifier string `json:"identifier"`
	ExpiryDate string `json:"expiryDate"`
	Token     string `json:"token"`
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

func (s *serviceManager) SplitProvisioningParameters(
	cpp map[string]interface{},
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
){
	pp := provisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil{
		return nil, nil, err
	}
	ppMap,err := service.GetMapFromStruct(pp)
	if err != nil {
		return nil, nil, err
	}
	return ppMap, nil, nil
}

func (s *serviceManager) SplitBindingParameters(
	params map[string]interface{},
) (
	service.BindingParameters,
	service.SecureBindingParameters,
	error,
) {
	return nil, nil, nil
}
