package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
			service.NewService(
				service.ServiceProperties{
					ID:          "58d9fbbd-7041-4dbe-aabe-6268cd31de84",
					Name:        "azure-cosmosdb-sql",
					Description: "Azure Cosmos DB (SQL API Database Account and Database)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API Database Account and Database)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
					},
				},
				m.sqlAllInOneManager,
				service.NewPlan(service.PlanProperties{
					ID:          "58d7223d-934e-4fb5-a046-0c67781eb24e",
					Name:        "sql-api",
					Description: "Azure CosmosDB With SQL API (Database Account and Database)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API Database Account and Database)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateProvisioningParamsSchema(), // nolint: lll
							UpdatingParametersSchema:     generateUpdatingParamsSchema(),     // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:             "6330de6f-a561-43ea-a15e-b99f44d183e6",
					Name:           "azure-cosmosdb-sql-account",
					Description:    "Azure Cosmos DB Database Account (SQL API)",
					ChildServiceID: "87c5132a-6d76-40c6-9621-0c7b7542571b",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API - Database Account Only)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
					},
				},
				m.sqlAccountManager,
				service.NewPlan(service.PlanProperties{
					ID:          "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name:        "account",
					Description: "Database Account with the SQL API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database Account Only)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateProvisioningParamsSchema(), // nolint: lll
							UpdatingParametersSchema:     generateUpdatingParamsSchema(),     // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "87c5132a-6d76-40c6-9621-0c7b7542571b",
					Name:        "azure-cosmosdb-sql-database",
					Description: "Azure Cosmos DB Database (SQL API - Database Only)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API - Database Only)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
					},
					ParentServiceID: "6330de6f-a561-43ea-a15e-b99f44d183e6",
				},
				m.sqlDatabaseManager,
				service.NewPlan(service.PlanProperties{
					ID:          "c821c68c-c8e0-4176-8cf2-f0ca582a07a3",
					Name:        "database",
					Description: "Azure CosmosDB (SQL API - Database only)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database only)",
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
					Name:        "azure-cosmosdb-mongo-account",
					Description: "Azure Cosmos DB Database Account (MongoDB API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"MongoDB",
					},
				},
				m.mongoAccountManager,
				service.NewPlan(service.PlanProperties{
					ID:          "86fdda05-78d7-4026-a443-1325928e7b02",
					Name:        "account",
					Description: "Database Account with the MongoDB API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateProvisioningParamsSchema(), // nolint: lll
							UpdatingParametersSchema:     generateUpdatingParamsSchema(),     // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "5f5252a0-6922-4a0c-a755-f9be70d7c79b",
					Name:        "azure-cosmosdb-graph-account",
					Description: "Azure Cosmos DB Database Account (Graph API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (Graph API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"Graph",
						"Gremlin",
					},
				},
				m.graphAccountManager,
				service.NewPlan(service.PlanProperties{
					ID:          "126a2c47-11a3-49b1-833a-21b563de6c04",
					Name:        "account",
					Description: "Database Account with the Graph API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Graph API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateProvisioningParamsSchema(), // nolint: lll
							UpdatingParametersSchema:     generateUpdatingParamsSchema(),     // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "37915cad-5259-470d-a7aa-207ba89ada8c",
					Name:        "azure-cosmosdb-table-account",
					Description: "Azure Cosmos DB Database Account (Table API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (Table API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"Table",
					},
				},
				m.tableAccountManager,
				service.NewPlan(service.PlanProperties{
					ID:          "c970b1e8-794f-4d7c-9458-d28423c08856",
					Name:        "account",
					Description: "Database Account with the Table API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Table API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateProvisioningParamsSchema(), // nolint: lll
							UpdatingParametersSchema:     generateUpdatingParamsSchema(),     // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "b8b56d60-4525-41d8-b3d8-8caa4dce0188",
					Name:        "azure-cosmosdb-sql-registered",
					Description: "Azure Cosmos DB From Registered (SQL API Database Account and Database)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From Registered (SQL API Database Account and Database)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
						"Disaster Recovery",
					},
				},
				m.sqlAllInOneRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "555ff2f7-336b-40f5-94aa-84d71d81d0af",
					Name:        "sql-api",
					Description: "Azure CosmosDB With SQL API (Database Account and Database)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API Database Account and Database)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateAllInOneRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "1b6e65c8-bab3-416b-8593-4dcad22e12d6",
					Name:        "azure-cosmosdb-sql-account-registered",
					Description: "Azure Cosmos DB Database Account From Registered (SQL API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From registered (SQL API - Database Account Only)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
						"Disaster Recovery",
					},
				},
				m.commonRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "49510891-89c5-4ba6-9207-e491c3f0ed90",
					Name:        "account",
					Description: "Database Account with the SQL API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database Account Only)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateCommonRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "beadd877-0c5a-476f-90ae-8922d17303e4",
					Name:        "azure-cosmosdb-sql-database-registered",
					Description: "Azure Cosmos DB Database From Registered (SQL API - Database Only)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From Registered (SQL API - Database Only)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
						"Disaster Recovery",
					},
				},
				m.sqlAllInOneRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "3d66321c-a38f-46d3-a24a-18619389a9b7",
					Name:        "database",
					Description: "Azure CosmosDB (SQL API - Database only)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database only)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateAllInOneRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "898c5053-8fa1-4ea9-89e3-301a4afc6bc2",
					Name:        "azure-cosmosdb-mongo-account-registered",
					Description: "Azure Cosmos DB Database Account From Registered (MongoDB API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From Registered (MongoDB API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"MongoDB",
						"Disaster Recovery",
					},
				},
				m.commonRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "468d259f-5fc3-4384-ab07-3e7efab4d587",
					Name:        "account",
					Description: "Database Account with the MongoDB API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateCommonRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "ab25a081-1a28-4132-96af-f60ef8201f75",
					Name:        "azure-cosmosdb-graph-account-registered",
					Description: "Azure Cosmos DB Database Account From Registered (Graph API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From Registered (Graph API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"Graph",
						"Gremlin",
						"Disaster Recovery",
					},
				},
				m.commonRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "d0f037ee-0abc-4f22-9667-4d27d773daab",
					Name:        "account",
					Description: "Database Account with the Graph API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Graph API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateCommonRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				service.ServiceProperties{
					ID:          "5f1c332b-c878-4e18-ba49-ccbfd2e46353",
					Name:        "azure-cosmosdb-table-account-registered",
					Description: "Azure Cosmos DB Database Account From Registered (Table API)",
					Metadata: service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB From Registered (Table API)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"Table",
						"Disaster Recovery",
					},
				},
				m.commonRegisteredManager,
				service.NewPlan(service.PlanProperties{
					ID:          "9c5a5bdd-a3d8-470b-9d98-957557a37a48",
					Name:        "account",
					Description: "Database Account with the Table API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Table API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: generateCommonRegisteredProvisioningParamsSchema(), // nolint: lll
						},
					},
				}),
			),
		}),
		nil
}
