package apimanagement

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"time"
)

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	pp := instance.ProvisioningParameters

	resourceGroup := pp.GetString("resourceGroup")
	apiName := pp.GetString("apiName")
	tenantClient := s.tenantAccessClient
	accessInformation, err := tenantClient.Get(context.TODO(), resourceGroup, apiName)
	if err != nil {
		return nil, err
	}

	identifier := *(accessInformation.ID)
	pk := *(accessInformation.PrimaryKey)
	expiry := time.Now().Add(time.Hour * 24 * 30).UTC()
	expiryDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.0000000Z", expiry.Year(), expiry.Month(), expiry.Day(), expiry.Hour(), expiry.Minute(), expiry.Second())
	key, err := generateKey(identifier, pk, expiryDate)
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("https://%s.management.azure-api.net/", apiName)
	return credentials{
		BaseURL:    baseURL,
		Identifier: identifier,
		ExpiryDate: expiryDate,
		Key:        key,
	}, nil
}

// This method is used to generate api management token, see
// https://docs.microsoft.com/en-us/rest/api/apimanagement/apimanagementrest/azure-api-management-rest-api-authentication
// for details
func generateKey(identifier string, key string, expiryDate string) (string, error) {
	toEncode := identifier + "\n" + expiryDate

	hashFunc := hmac.New(sha512.New, []byte(key))
	if _, err := hashFunc.Write([]byte(toEncode)); err != nil {
		return "", err
	}
	encryptedBytes := hashFunc.Sum(nil)
	token := base64.StdEncoding.EncodeToString(encryptedBytes)
	return token, nil
}
