package cosmosdb

import (
	"context"
	"fmt"

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
	dt.DatabaseName = instance.ProvisioningParameters.GetString("databaseName")
	return dt, nil
}
