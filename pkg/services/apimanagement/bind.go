package apimanagement

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"context"
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
)

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}

	resourceGroup := instance.ResourceGroup
	apiName := dt.apiName
	tenantClient := s.tenantAccessClient
	accessInformation, err := tenantClient.Get(context.TODO(), resourceGroup, apiName)
	if err != nil {
		return nil, err
	}

	identifier := *(accessInformation.ID)
	key := *(accessInformation.PrimaryKey)
	expiry := time.Now().Add(time.Hour * 24 * 30).UTC()
	token,err := generateToken(identifier, key, expiry)
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("https://%s.management.azure-api.net/", apiName)
	return credentials {
		BaseURL: baseURL,
		Token : token,
	}, nil
}

//This method is used to generate api management token, see
// https://docs.microsoft.com/en-us/rest/api/apimanagement/apimanagementrest/azure-api-management-rest-api-authentication
// for details
func generateToken(identifier string, key string, expiry time.Time) (string, error) {
	expiryDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.0000000Z",expiry.Year(),expiry.Month(),expiry.Day(),expiry.Hour(),expiry.Minute(),expiry.Second())
	toEncode := identifier + "\n" + expiryDate

	hashFunc := hmac.New(sha512.New, []byte(key))
	if _, err := hashFunc.Write([]byte(toEncode)); err != nil {
		return "", err
	}
	encryptedBytes := hashFunc.Sum(nil)
	token := base64.StdEncoding.EncodeToString(encryptedBytes)
	return token, nil
}
