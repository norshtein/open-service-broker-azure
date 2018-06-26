package apimanagement

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "157d6551-c63a-4002-bc03-87c248ad42a1",
				Name:        "azure-api-management",
				Description: "Azure API Management (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure API Management",
					ImageURL: "https://azure.microsoft.com/svghandler/api-management/" +
						"?width=200",
					LongDescription: "offers a scalable API gateway for securing, publishing, and " +
						"analyzing APIs and microservices to internal and external consumers (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/api-management/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "API", "Management"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "df428328-8b77-473e-8bef-3589ef8e612f",
				Name:        "developer",
				Description: "Developer Tier(No SLA)",
				Free:        false,
				Extended: map[string]interface{}{
					"tier": "Developer",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Developer Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "287250d0-f10a-42ea-93ab-53f8bb454b14",
				Name:        "basic",
				Description: "Basic Tier",
				Free:        false,
				Extended: map[string]interface{}{
					"tier": "Basic",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic Tier(99.9 SLA, %)",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "f0392d05-71c7-45c1-bd60-cca4ded8dc8a",
				Name:        "standard",
				Description: "Standard Tier(99.9 SLA, %)",
				Free:        false,
				Extended: map[string]interface{}{
					"tier": "Standard",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "f171ec17-d247-404e-b5b8-2c638b0bc59c",
				Name:        "premium",
				Description: "Premium Tier(99.95* SLA, %)",
				Free:        false,
				Extended: map[string]interface{}{
					"tier": "Premium",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Premium Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
		),
	}), nil
}
