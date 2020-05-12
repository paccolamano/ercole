// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"context"
	"testing"

	"github.com/amreo/ercole-services/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchHosts() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, "", false, -1, -1, "", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, "", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, "_id", true, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"CPUCores":       1,
						"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
						"CPUThreads":     2,
						"Cluster":        "Puzzait",
						"CreatedAt":      "2020-05-04T16:09:46.608+02:00",
						"Databases":      "",
						"Environment":    "PROD",
						"HostType":       "virtualization",
						"Hostname":       "test-virt",
						"Kernel":         "3.10.0-862.9.1.el7.x86_64",
						"Location":       "Italy",
						"MemTotal":       3,
						"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
						"OracleCluster":  false,
						"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
						"Socket":         2,
						"SunCluster":     false,
						"SwapTotal":      4,
						"Type":           "VMWARE",
						"VeritasCluster": false,
						"Version":        "1.6.1",
						"Virtual":        true,
						"_id":            "5eb0222a45d85f4193704944",
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          false,
					"Number":        0,
					"Size":          1,
					"TotalElements": 3,
					"TotalPages":    3,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, "CreatedAt", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":     2,
				"Cluster":        "Puzzait",
				"CreatedAt":      "2020-05-04T16:09:46.608+02:00",
				"Databases":      "",
				"Environment":    "PROD",
				"HostType":       "virtualization",
				"Hostname":       "test-virt",
				"Kernel":         "3.10.0-862.9.1.el7.x86_64",
				"Location":       "Italy",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      4,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "1.6.1",
				"Virtual":        true,
				"_id":            "5eb0222a45d85f4193704944",
			},
			{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":     2,
				"Cluster":        interface{}(nil),
				"CreatedAt":      "2020-04-24T13:50:05.46+02:00",
				"Databases":      "ERCOLE",
				"Environment":    "TST",
				"HostType":       "oracledb",
				"Hostname":       "test-small",
				"Kernel":         "3.10.0-514.el7.x86_64",
				"Location":       "Germany",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   interface{}(nil),
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      1,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "latest",
				"Virtual":        true,
				"_id":            "5ea2d26d20d55cbdc35022b4",
			},
			{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":     2,
				"Cluster":        "Puzzait",
				"CreatedAt":      "2020-04-15T08:46:58.471+02:00",
				"Databases":      "ERCOLE",
				"Environment":    "TST",
				"HostType":       "oracledb",
				"Hostname":       "test-db",
				"Kernel":         "3.10.0-514.el7.x86_64",
				"Location":       "Germany",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      1,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "latest",
				"Virtual":        true,
				"_id":            "5e96ade270c184faca93fe36",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search1", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search2", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"test-db", "ERCOLE"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":     2,
				"Cluster":        "Puzzait",
				"CreatedAt":      "2020-04-15T08:46:58.471+02:00",
				"Databases":      "ERCOLE",
				"Environment":    "TST",
				"HostType":       "oracledb",
				"Hostname":       "test-db",
				"Kernel":         "3.10.0-514.el7.x86_64",
				"Location":       "Germany",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      1,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "latest",
				"Virtual":        true,
				"_id":            "5e96ade270c184faca93fe36",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search3", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"Puzzait"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":     2,
				"Cluster":        "Puzzait",
				"CreatedAt":      "2020-05-04T16:09:46.608+02:00",
				"Databases":      "",
				"Environment":    "PROD",
				"HostType":       "virtualization",
				"Hostname":       "test-virt",
				"Kernel":         "3.10.0-862.9.1.el7.x86_64",
				"Location":       "Italy",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      4,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "1.6.1",
				"Virtual":        true,
				"_id":            "5eb0222a45d85f4193704944",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("lms_mode", func(t *testing.T) {
		out, err := m.db.SearchHosts("lms", []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"ConnectString":            "",
				"CoresPerProcessor":        1,
				"DBInstanceName":           "ERCOLE",
				"Environment":              "TST",
				"Features":                 "Diagnostics Pack",
				"Notes":                    "",
				"OperatingSystem":          "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
				"PhysicalCores":            2,
				"PhysicalServerName":       "Puzzait",
				"PluggableDatabaseName":    "",
				"ProcessorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"ProcessorSpeed":           "2.53GHz",
				"Processors":               2,
				"ProductEdition":           "Enterprise",
				"ProductVersion":           "12",
				"RacNodeNames":             "",
				"ServerPurchaseDate":       "",
				"ThreadsPerCore":           2,
				"VirtualServerName":        "test-db",
				"VirtualizationTechnology": "VMWARE",
				"_id":                      utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}