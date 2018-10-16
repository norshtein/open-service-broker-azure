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

	databaseAccountName := instance.ProvisioningParameters.GetString("accountName") // nolint: lll
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup") // nolint: lll

	databaseAccount, err := s.databaseAccountsClient.Get(
		ctx,
		resourceGroupName,
		databaseAccountName,
	)
	if err != nil {
		return nil, err
	}
	fqdn := *(databaseAccount.DatabaseAccountProperties.DocumentEndpoint)

	keys, err := s.databaseAccountsClient.ListKeys(
		ctx,
		resourceGroupName,
		databaseAccountName,
	)
	if err != nil {
		return nil, err
	}
	pk := *(keys.PrimaryMasterKey)

	dt := &sqlAllInOneInstanceDetails{}
	dt.DatabaseAccountName = databaseAccountName
	dt.PrimaryKey = service.SecureString(pk)
	dt.FullyQualifiedDomainName = fqdn
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
