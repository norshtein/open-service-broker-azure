// +build !unit

package lifecycle

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var cosmosdbTestCases = []serviceLifecycleTestCase{
	{ // SQL API
		group:     "cosmosdb",
		name:      "sql-api-account-only",
		serviceID: "6330de6f-a561-43ea-a15e-b99f44d183e6",
		planID:    "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
		provisioningParameters: map[string]interface{}{
			"alias":    "cosmos-account",
			"location": "eastus",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []interface{}{"0.0.0.0/0"},
			},
			"consistencyPolicy": map[string]interface{}{
				"defaultConsistencyLevel": "Session",
			},
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
		updatingParameters: map[string]interface{}{
			"readRegions": []interface{}{"centralus"},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:     "cosmosdb",
				name:      "database-only",
				serviceID: "87c5132a-6d76-40c6-9621-0c7b7542571b",
				planID:    "c821c68c-c8e0-4176-8cf2-f0ca582a07a3",
				provisioningParameters: map[string]interface{}{
					"parentAlias": "cosmos-account",
				},
			},
		},
	},
	{ // Table API
		group:     "cosmosdb",
		name:      "table-api-account-only",
		serviceID: "37915cad-5259-470d-a7aa-207ba89ada8c",
		planID:    "c970b1e8-794f-4d7c-9458-d28423c08856",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []interface{}{"0.0.0.0/0"},
			},
			"consistencyPolicy": map[string]interface{}{
				"defaultConsistencyLevel": "Session",
			},
			"readRegions": []interface{}{"eastus2"},
		},
		updatingParameters: map[string]interface{}{
			"readRegions": []interface{}{"centralus"},
		},
	},
}

func testMongoDBCreds(credentials map[string]interface{}) error {
	// The following process is based on
	// https://docs.microsoft.com/en-us/azure/cosmos-db/create-mongodb-golang

	// DialInfo holds options for establishing a session with a MongoDB cluster.
	dialInfo := &mgo.DialInfo{
		Addrs: []string{
			fmt.Sprintf(
				"%s:%d",
				credentials["host"].(string),
				int(credentials["port"].(float64)),
			),
		},
		Timeout:  60 * time.Second,
		Database: "database",
		Username: credentials["username"].(string),
		Password: credentials["password"].(string),
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}

	// Create a session which maintains a pool of socket connections
	// to our Azure Cosmos DB MongoDB database.
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetSafe(&mgo.Safe{})

	collection := session.DB("database").C("package")

	// Model
	type Package struct {
		ID            bson.ObjectId `bson:"_id,omitempty"`
		FullName      string
		Description   string
		StarsCount    int
		ForksCount    int
		LastUpdatedBy string
	}

	// insert Document in collection
	err = collection.Insert(&Package{
		FullName:      "react",
		Description:   "A framework for building native apps with React.",
		ForksCount:    11392,
		StarsCount:    48794,
		LastUpdatedBy: "shergin",
	})

	return err
}
