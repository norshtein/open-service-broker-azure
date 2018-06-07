package apimanagement

// nolint: lll
var armTemplateBytes = []byte(`
{
		"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": {
			"name": {
				"type": "String"
			},
			"location": {
				"type": "String"
			},
			"adminEmail": {
				"type": "String"
			},
			"orgName": {
				"type": "String"
			},
			"tier": {
				"type": "String"
			}
		},
		"resources": [
			{
				"type": "Microsoft.ApiManagement/service",
				"sku": {
					"name": "[parameters('tier')]",
					"capacity": 1
				},
				"name": "[parameters('name')]",
				"apiVersion": "2018-01-01",
				"location": "[parameters('location')]",
				"properties": {
					"publisherEmail": "[parameters('adminEmail')]",
					"publisherName": "[parameters('orgName')]"
				}
			}
		]
	}
`)

