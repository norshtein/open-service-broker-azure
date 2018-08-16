package cosmosdb

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("fillInInstanceDetails", s.fillInInstanceDetails), // nolint:lll
	)
}

func (s *sqlAllInOneRegisteredManager) fillInInstanceDetails(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pk, err := getPrimaryKey(
		ctx,
		instance.ProvisioningParameters,
		s.databaseAccountsClient,
	)
	if err != nil {
		return nil, err
	}

	dt := &sqlAllInOneInstanceDetails{}
	dt.PrimaryKey = pk
	dt.FullyQualifiedDomainName = fmt.Sprintf(
		"https://%s.documents.azure.com:443/",
		instance.ProvisioningParameters.GetString("accountName"),
	)
	dt.ConnectionString = service.SecureString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		),
	)
	if err = validateDatabase(
		instance.ProvisioningParameters.GetString("accountName"),
		instance.ProvisioningParameters.GetString("databaseName"),
		string(dt.PrimaryKey),
	); err != nil {
		return nil, err
	}
	dt.DatabaseName = instance.ProvisioningParameters.GetString("databaseName")
	return dt, nil
}

func validateDatabase(
	accountName string,
	databaseName string,
	key string,
) error {
	req, err := createRequest(
		accountName,
		"GET",
		databaseName,
		key,
		nil,
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
	if resp.StatusCode != 200 { // CosmosDB returns a 200 on success
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"error validating database %d : unable to get body",
				resp.StatusCode,
			)
		}
		return fmt.Errorf(
			"error validating database %d : %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}
