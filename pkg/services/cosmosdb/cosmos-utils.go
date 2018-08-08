package cosmosdb

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/tidwall/gjson"
)

// This method implements the CosmosDB API authentication token generation
// scheme. For reference, please see the CosmosDB REST API at:
// https://aka.ms/Fyra7j
func generateAuthToken(verb, id, date, key string) (string, error) {
	resource := "dbs"
	var resourceID string
	if id != "" {
		resourceID = fmt.Sprintf("%s/%s", strings.ToLower(resource), id)
	} else {
		resourceID = id
	}
	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		strings.ToLower(verb),
		strings.ToLower(resource),
		resourceID,
		strings.ToLower(date),
		"",
	)

	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	hmac := hmac.New(sha256.New, decodedKey)
	_, err := hmac.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	b := hmac.Sum(nil)
	authHash := base64.StdEncoding.EncodeToString(b)
	authHeader := url.QueryEscape("type=master&ver=1.0&sig=" + authHash)
	return authHeader, nil
}

func createRequest(
	accountName string,
	method string,
	resourceID string,
	key string,
	body interface{},
) (*http.Request, error) {
	resourceType := "dbs" // If we support other types, parameterize this
	path := fmt.Sprintf("%s/%s", resourceType, resourceID)
	url := fmt.Sprintf("https://%s.documents.azure.com/%s", accountName, path)
	var buf *bytes.Buffer
	var err error
	var req *http.Request
	if body != nil {
		var b []byte
		b, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf(
				"error building comsosdb request body: %s",
				err,
			)
		}
		buf = bytes.NewBuffer(b)
		req, err = http.NewRequest(method, url, buf)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("error building comsosdb request: %s", err)
	}

	dateStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	authHeader, err := generateAuthToken(
		method,
		resourceID,
		dateStr,
		key,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Ms-Date", dateStr)
	req.Header.Add("X-Ms-version", "2017-02-22")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

func createDatabase(
	accountName string,
	id string,
	key string,
) error {
	request := &databaseCreationRequest{
		ID: id,
	}
	databaseName := ""
	req, err := createRequest(
		accountName,
		"POST",
		databaseName,
		key,
		request,
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"error making create comsosdb database request: %s",
			err,
		)
	}
	if resp.StatusCode != 201 { // CosmosDB returns a 201 on success
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"error creating database %d : unable to get body",
				resp.StatusCode,
			)
		}
		return fmt.Errorf(
			"error creating database %d : %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}

func deleteDatabase(
	databaseAccount string,
	databaseName string,
	key string,
) error {
	req, err := createRequest(
		databaseAccount,
		"DELETE",
		databaseName,
		key,
		nil, //No Body here
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"error making delete comsosdb database request: %s",
			err,
		)
	}
	if resp.StatusCode != 204 { // CosmosDB returns a 204 on success
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"error deleting database %d : unable to get body",
				resp.StatusCode,
			)
		}
		return fmt.Errorf(
			"error deleting database %d : %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}

// The deployment will return success once the write region is created, ignoring the status of read regions, so we must implement detection logic by ourselves.
// For now, this method will return on either context is cancelled or every region's state is "succeeded" in seven consecutive check.
// The reason why we need seven consecutive check is that the read region is created one by one, there is a small gap between
// the finishment of previous creation and the start of the next creation. By this check, we can detect gaps shorter than 1 mintue,
// and report success within 70 seconds after completion.
func (c *cosmosAccountManager) waitForReadRegionsReady(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Bug here for sql-all-in-one
	dt := instance.Details.(*cosmosdbInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := c.databaseAccountsClient

	ticker := time.NewTicker(time.Second * 10)
	previousSucceededTimes := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			result, err := databaseAccountClient.Get(childCtx, resourceGroupName, accountName)
			if err != nil {
				return nil, err
			}
			resultJSONBytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			resultJSONString := string(resultJSONBytes)

			//Check whether every read location's state is "Succeeded"
			allSucceed := true
			readLocations := gjson.Get(resultJSONString, "properties.readLocations")
			readLocations.ForEach(func(key, value gjson.Result) bool {
				state := value.Get("provisioningState").String()
				if state != "Succeeded" {
					allSucceed = false
					previousSucceededTimes = 0
					return false
				}
				return true
			})
			if allSucceed && previousSucceededTimes >= 7 {
				return dt, nil
			} else if allSucceed {
				previousSucceededTimes++
			}
		}
	}
}

// For sqlAllInOneManager, the real type of `instance.Details` is `*sqlAllInOneInstanceDetails`,
// so type assertion must be changed. Expect line 249, this function is totally the same as previous one.
// Do you have any good idea make the code cleaner?
func (s *sqlAllInOneManager) waitForReadRegionsReady(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := s.databaseAccountsClient

	ticker := time.NewTicker(time.Second * 10)
	previousSucceededTimes := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			result, err := databaseAccountClient.Get(childCtx, resourceGroupName, accountName)
			if err != nil {
				return nil, err
			}
			resultJSONBytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			resultJSONString := string(resultJSONBytes)

			//Check whether every read location's state is "Succeeded"
			allSucceed := true
			readLocations := gjson.Get(resultJSONString, "properties.readLocations")
			readLocations.ForEach(func(key, value gjson.Result) bool {
				state := value.Get("provisioningState").String()
				if state != "Succeeded" {
					allSucceed = false
					previousSucceededTimes = 0
					return false
				}
				return true
			})
			if allSucceed && previousSucceededTimes >= 7 {
				return dt, nil
			} else if allSucceed {
				previousSucceededTimes++
			}
		}
	}
}

func validateReadRegions(
	context string,
	regions []string,
) error {
	for i := range regions {
		region := regions[i]
		if !allowedReadRegions[region] {
			return service.NewValidationError(
				fmt.Sprintf("%s.allowedReadRegion", context),
				fmt.Sprintf("given region %s is not allowed", region),
			)
		}
	}
	return nil
}

// Allowed CosmosDB read regions
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
	"brazilsouth":        true,
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
