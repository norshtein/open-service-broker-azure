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

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
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

// Although there is a `waitForCompletion` method, but that method is implemented by checking HTTP status code,
// but unfortunately, in cosmosSDK, the REST API behind `createOrUpdate` will always return "200 OK",
// so `waitForCompletion` method cannot detect whether the update has finished, we must implement detection logic by ourselves.
// For now, this method will return on either context is cancelled or every region's state is "succeeded"
func waitForRegionCreationCompletion(ctx context.Context, dac cosmosSDK.DatabaseAccountsClient, resourceGroupName string, accountName string) error {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			result, err := dac.Get(childCtx, resourceGroupName, accountName)
			if err != nil {
				return err
			}
			resultJSONBytes, err := json.Marshal(result)
			if err != nil {
				return err
			}
			resultJSONString := string(resultJSONBytes)

			//Check whether every read location's state is "Succeeded"
			allSucceed := true
			readLocations := gjson.Get(resultJSONString, "properties.readLocations")
			readLocations.ForEach(func(key, value gjson.Result) bool {
				state := value.Get("provisioningState").String()
				if state != "Succeeded" {
					fmt.Printf("State is %s\n", state)
					allSucceed = false
					return false
				}
				return true
			})
			if allSucceed {
				return nil
			}
		}
	}
}

func contructLocation(accountName, locationName string, failoverPriority int32) cosmosSDK.Location {
	return cosmosSDK.Location{
		FailoverPriority: &failoverPriority,
		LocationName:     &locationName,
	}
}
