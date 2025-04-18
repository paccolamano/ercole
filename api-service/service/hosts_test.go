// Copyright (c) 2022 Sorint.lab S.p.A.
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

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchHosts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Angola-1dac9f7418db9b52c259ce4ba087cdb6",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.844+02:00"),
			"Databases":                     "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":                   "PROD",
			"Hostname":                      "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "suspended-290dce22a939f3868f8f23a6e1f57dd8",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Puzzait",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.869+02:00"),
			"Databases":                     "",
			"Environment":                   "PROD",
			"Hostname":                      "test-virt",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	db.EXPECT().SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchHosts_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	)
	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchHostsAsLMS(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	hosts := []map[string]interface{}{
		{
			"coresPerProcessor":        1,
			"dbInstanceName":           "ERCOLE",
			"environment":              "TST",
			"licenseMetricAllocated":   "processor",
			"operatingSystem":          "Red Hat Enterprise Linux",
			"options":                  "Diagnostics Pack",
			"physicalCores":            2,
			"physicalServerName":       "erclin7dbx",
			"pluggableDatabaseName":    "",
			"processorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"processorSpeed":           "2.53GHz",
			"processors":               2,
			"productLicenseAllocated":  "EE",
			"productVersion":           "12",
			"threadsPerCore":           2,
			"usedManagementPacks":      "Diagnostics Pack",
			"usingLicenseCount":        0.5,
			"virtualServerName":        "itl-csllab-112.sorint.localpippo",
			"virtualizationTechnology": "VMware",
			"_id":                      utils.Str2oid("5efc38ab79f92e4cbf283b03"),
			"createdAt":                utils.PDT("2021-12-05T00:00:00+02:00"),
		},
		{
			"coresPerProcessor":        4,
			"dbInstanceName":           "rudeboy-fb3160a04ffea22b55555bbb58137f77",
			"environment":              "SVIL",
			"licenseMetricAllocated":   "processor",
			"operatingSystem":          "Red Hat Enterprise Linux",
			"options":                  "",
			"physicalCores":            8,
			"physicalServerName":       "",
			"pluggableDatabaseName":    "",
			"processorModel":           "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz",
			"processorSpeed":           "2.93GHz",
			"processors":               2,
			"productLicenseAllocated":  "EE",
			"productVersion":           "11",
			"threadsPerCore":           2,
			"usedManagementPacks":      "",
			"usingLicenseCount":        4.0,
			"virtualServerName":        "publicitate-36d06ca83eafa454423d2097f4965517",
			"virtualizationTechnology": "",
			"_id":                      utils.Str2oid("5efc38ab79f92e4cbf283b04"),
			"createdAt":                utils.PDT("2020-12-05T00:00:00+02:00"),
		},
		{
			"coresPerProcessor":        2,
			"dbInstanceName":           "buu",
			"environment":              "SVIL",
			"licenseMetricAllocated":   "processor",
			"operatingSystem":          "Red Hat Enterprise Linux",
			"options":                  "",
			"physicalCores":            4,
			"physicalServerName":       "",
			"pluggableDatabaseName":    "",
			"processorModel":           "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz",
			"processorSpeed":           "2.93GHz",
			"processors":               2,
			"productLicenseAllocated":  "EE",
			"productVersion":           "11",
			"threadsPerCore":           2,
			"usedManagementPacks":      "",
			"usingLicenseCount":        4.0,
			"virtualServerName":        "itl-csllab-112.sorint.localbuu",
			"virtualizationTechnology": "",
			"_id":                      utils.Str2oid("5efc38ab79f92e4cbf283b05"),
			"createdAt":                utils.PDT("2020-12-05T00:00:00+02:00"),
			"dismissedAt":              utils.PDT("2021-05-10T00:00:00+02:00"),
		},
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}

	filterslms := dto.SearchHostsAsLMS{
		SearchHostsFilters: filters,
		From:               utils.P("2020-06-10T11:54:59Z"),
		To:                 utils.P("2021-06-10T11:54:59Z"),
	}

	t.Run("with no contracts", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().
				SearchHosts("lms", gomock.Any()).
				DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
					assert.EqualValues(t, filterslms.SearchHostsFilters, actual)

					return hosts, nil
				}),
			db.EXPECT().
				ListOracleDatabaseContracts(gomock.Any()).
				Return([]dto.OracleDatabaseContractFE{}, nil),
			db.EXPECT().
				GetListValidHostsByRangeDates(filterslms.From, filterslms.To).
				DoAndReturn(func(from time.Time, to time.Time) ([]string, error) {
					return []string{}, nil
				}).Times(1),
			db.EXPECT().
				GetListDismissedHostsByRangeDates(filterslms.From, filterslms.To).
				DoAndReturn(func(from time.Time, to time.Time) ([]string, error) {
					return []string{}, nil
				}).Times(1),
			db.EXPECT().
				GetListValidHostsByRangeDates(filterslms.From, filterslms.To).
				DoAndReturn(func(from time.Time, to time.Time) ([]string, error) {
					return []string{}, nil
				}).Times(1),
			db.EXPECT().
				GetListDismissedHostsByRangeDates(filterslms.From, filterslms.To).
				DoAndReturn(func(from time.Time, to time.Time) ([]string, error) {
					return []string{}, nil
				}).Times(1),
		)

		sp, err := as.SearchHostsAsLMS(filterslms)
		assert.NoError(t, err)

		assert.Equal(t, "erclin7dbx", sp.GetCellValue("Database_&_EBS_DB_Tier", "B4"))
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", sp.GetCellValue("Database_&_EBS_DB_Tier", "C4"))
		assert.Equal(t, "VMware", sp.GetCellValue("Database_&_EBS_DB_Tier", "D4"))
		assert.Equal(t, "ERCOLE", sp.GetCellValue("Database_&_EBS_DB_Tier", "E4"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "F4"))
		assert.Equal(t, "TST", sp.GetCellValue("Database_&_EBS_DB_Tier", "G4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue("Database_&_EBS_DB_Tier", "H4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue("Database_&_EBS_DB_Tier", "I4"))
		assert.Equal(t, "12", sp.GetCellValue("Database_&_EBS_DB_Tier", "N4"))
		assert.Equal(t, "EE", sp.GetCellValue("Database_&_EBS_DB_Tier", "O4"))
		assert.Equal(t, "processor", sp.GetCellValue("Database_&_EBS_DB_Tier", "P4"))
		assert.Equal(t, "0.5", sp.GetCellValue("Database_&_EBS_DB_Tier", "Q4"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AC4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AD4"))
		assert.Equal(t, "1", sp.GetCellValue("Database_&_EBS_DB_Tier", "AE4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AF4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AG4"))
		assert.Equal(t, "2.53GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AH4"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue("Database_&_EBS_DB_Tier", "AJ4"))

		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "B5"))
		assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sp.GetCellValue("Database_&_EBS_DB_Tier", "C5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "D5"))
		assert.Equal(t, "rudeboy-fb3160a04ffea22b55555bbb58137f77", sp.GetCellValue("Database_&_EBS_DB_Tier", "E5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "F5"))
		assert.Equal(t, "SVIL", sp.GetCellValue("Database_&_EBS_DB_Tier", "G5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "H5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "I5"))
		assert.Equal(t, "11", sp.GetCellValue("Database_&_EBS_DB_Tier", "N5"))
		assert.Equal(t, "EE", sp.GetCellValue("Database_&_EBS_DB_Tier", "O5"))
		assert.Equal(t, "processor", sp.GetCellValue("Database_&_EBS_DB_Tier", "P5"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "Q5"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AC5"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AD5"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "AE5"))
		assert.Equal(t, "8", sp.GetCellValue("Database_&_EBS_DB_Tier", "AF5"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AG5"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AH5"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue("Database_&_EBS_DB_Tier", "AJ5"))

		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "B6"))
		assert.Equal(t, "itl-csllab-112.sorint.localbuu", sp.GetCellValue("Database_&_EBS_DB_Tier", "C6"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "D6"))
		assert.Equal(t, "buu", sp.GetCellValue("Database_&_EBS_DB_Tier", "E6"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "F6"))
		assert.Equal(t, "SVIL", sp.GetCellValue("Database_&_EBS_DB_Tier", "G6"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "H6"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "I6"))
		assert.Equal(t, "11", sp.GetCellValue("Database_&_EBS_DB_Tier", "N6"))
		assert.Equal(t, "EE", sp.GetCellValue("Database_&_EBS_DB_Tier", "O6"))
		assert.Equal(t, "processor", sp.GetCellValue("Database_&_EBS_DB_Tier", "P6"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "Q6"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AC6"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AD6"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AE6"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "AF6"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AG6"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AH6"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue("Database_&_EBS_DB_Tier", "AJ6"))
	})

	t.Run("with contracts", func(t *testing.T) {
		contracts := []dto.OracleDatabaseContractFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaa"),
				ContractID:               "",
				CSI:                      "csi001",
				LicenseTypeID:            "",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "publicitate-36d06ca83eafa454423d2097f4965517"}},
				LicensesPerCore:          0,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 0,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaa"),
				ContractID:               "",
				CSI:                      "csi002",
				LicenseTypeID:            "",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "publicitate-36d06ca83eafa454423d2097f4965517"}},
				LicensesPerCore:          0,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 0,
				AvailableLicensesPerUser: 0,
			},
		}

		gomock.InOrder(
			db.EXPECT().
				SearchHosts("lms", gomock.Any()).
				DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
					assert.EqualValues(t, filterslms.SearchHostsFilters, actual)

					return hosts, nil
				}),
			db.EXPECT().
				ListOracleDatabaseContracts(gomock.Any()).
				Return(contracts, nil),
		)

		sp, err := as.SearchHostsAsLMS(filterslms)
		assert.NoError(t, err)

		sheet := "Database_&_EBS_DB_Tier"
		assert.Equal(t, "erclin7dbx", sp.GetCellValue(sheet, "B4"))
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", sp.GetCellValue(sheet, "C4"))
		assert.Equal(t, "VMware", sp.GetCellValue(sheet, "D4"))
		assert.Equal(t, "ERCOLE", sp.GetCellValue(sheet, "E4"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "F4"))
		assert.Equal(t, "TST", sp.GetCellValue(sheet, "G4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue(sheet, "H4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue(sheet, "I4"))
		assert.Equal(t, "12", sp.GetCellValue(sheet, "N4"))
		assert.Equal(t, "EE", sp.GetCellValue(sheet, "O4"))
		assert.Equal(t, "processor", sp.GetCellValue(sheet, "P4"))
		assert.Equal(t, "0.5", sp.GetCellValue(sheet, "Q4"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "R4"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sp.GetCellValue(sheet, "AC4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AD4"))
		assert.Equal(t, "1", sp.GetCellValue(sheet, "AE4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AF4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AG4"))
		assert.Equal(t, "2.53GHz", sp.GetCellValue(sheet, "AH4"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue(sheet, "AJ4"))

		assert.Equal(t, "", sp.GetCellValue(sheet, "B5"))
		assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sp.GetCellValue(sheet, "C5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "D5"))
		assert.Equal(t, "rudeboy-fb3160a04ffea22b55555bbb58137f77", sp.GetCellValue(sheet, "E5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "F5"))
		assert.Equal(t, "SVIL", sp.GetCellValue(sheet, "G5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "H5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "I5"))
		assert.Equal(t, "11", sp.GetCellValue(sheet, "N5"))
		assert.Equal(t, "EE", sp.GetCellValue(sheet, "O5"))
		assert.Equal(t, "processor", sp.GetCellValue(sheet, "P5"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "Q5"))
		assert.Equal(t, "csi001, csi002", sp.GetCellValue(sheet, "R5"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue(sheet, "AC5"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AD5"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "AE5"))
		assert.Equal(t, "8", sp.GetCellValue(sheet, "AF5"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AG5"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue(sheet, "AH5"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue(sheet, "AJ5"))

		assert.Equal(t, "", sp.GetCellValue(sheet, "B6"))
		assert.Equal(t, "itl-csllab-112.sorint.localbuu", sp.GetCellValue(sheet, "C6"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "D6"))
		assert.Equal(t, "buu", sp.GetCellValue(sheet, "E6"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "F6"))
		assert.Equal(t, "SVIL", sp.GetCellValue(sheet, "G6"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "H6"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "I6"))
		assert.Equal(t, "11", sp.GetCellValue(sheet, "N6"))
		assert.Equal(t, "EE", sp.GetCellValue(sheet, "O6"))
		assert.Equal(t, "processor", sp.GetCellValue(sheet, "P6"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "Q6"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "R6"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue(sheet, "AC6"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AD6"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AE6"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "AF6"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AG6"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue(sheet, "AH6"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue(sheet, "AJ6"))
	})
}

func TestSearchHostsAsXLSX(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	expectedRes := []dto.HostDataSummary{
		{
			ID:           "5efc38ab79f92e4cbf283b0b",
			CreatedAt:    utils.P("2020-07-01T09:18:03.715+02:00"),
			Hostname:     "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			Location:     "Italy",
			Environment:  "PROD",
			AgentVersion: "latest",
			Info: model.Host{
				Hostname:                      "",
				CPUModel:                      "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz",
				CPUFrequency:                  "",
				CPUSockets:                    1,
				CPUCores:                      24,
				CPUThreads:                    48,
				ThreadsPerCore:                0,
				CoresPerSocket:                0,
				HardwareAbstraction:           "PH",
				HardwareAbstractionTechnology: "PH",
				Kernel:                        "Linux 4.1.12-124.26.12.el7uek.x86_64",
				KernelVersion:                 "",
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.6",
				MemoryTotal:                   376,
				SwapTotal:                     23,
			},
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       true,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			VirtualizationNode: "",
			Cluster:            "",
			Databases:          map[string][]string{},
			MissingDatabases:   []model.MissingDatabase{},
		},

		{
			ID:           "5efc38ab79f92e4cbf283b13",
			CreatedAt:    utils.P("2020-07-01T09:18:03.726+02:00"),
			Hostname:     "test-db",
			Location:     "Germany",
			Environment:  "TST",
			AgentVersion: "latest",
			Info: model.Host{
				Hostname:                      "",
				CPUModel:                      "Intel(R) Xeon(R) CPU E5630  @ 2.53GHz",
				CPUFrequency:                  "",
				CPUSockets:                    2,
				CPUCores:                      1,
				CPUThreads:                    2,
				ThreadsPerCore:                0,
				CoresPerSocket:                0,
				HardwareAbstraction:           "VIRT",
				HardwareAbstractionTechnology: "VMWARE",
				Kernel:                        "Linux 3.10.0-514.el7.x86_64",
				KernelVersion:                 "",
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.6",
				MemoryTotal:                   3,
				SwapTotal:                     1,
			},
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
			Cluster:            "Puzzait",
			Databases:          map[string][]string{},
		},
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}

	t.Run("Success", func(t *testing.T) {
		db.EXPECT().
			GetHostDataSummaries(filters).
			Return(expectedRes, nil)

		sp, err := as.SearchHostsAsXLSX(filters)
		assert.NoError(t, err)

		assert.Equal(t, "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", sp.GetCellValue("Hosts", "A2"))
		assert.Equal(t, "Bare metal", sp.GetCellValue("Hosts", "B2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "C2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "D2"))
		assert.Equal(t, "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz", sp.GetCellValue("Hosts", "E2"))
		assert.Equal(t, "48", sp.GetCellValue("Hosts", "F2"))
		assert.Equal(t, "24", sp.GetCellValue("Hosts", "G2"))
		assert.Equal(t, "1", sp.GetCellValue("Hosts", "H2"))
		assert.Equal(t, "latest", sp.GetCellValue("Hosts", "I2"))
		assert.Equal(t, "7/1/20 09:18", sp.GetCellValue("Hosts", "J2"))
		assert.Equal(t, "PROD", sp.GetCellValue("Hosts", "K2"))
		assert.Equal(t, "Italy", sp.GetCellValue("Hosts", "L2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "M2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "N2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "O2"))
		assert.Equal(t, "Red Hat Enterprise Linux - 7.6", sp.GetCellValue("Hosts", "P2"))
		assert.Equal(t, "1", sp.GetCellValue("Hosts", "Q2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "R2"))
		assert.Equal(t, "376", sp.GetCellValue("Hosts", "S2"))
		assert.Equal(t, "23", sp.GetCellValue("Hosts", "T2"))
		assert.Equal(t, "ClusterWare", sp.GetCellValue("Hosts", "U2"))
		assert.Equal(t, "0", sp.GetCellValue("Hosts", "V2"))

		assert.Equal(t, "test-db", sp.GetCellValue("Hosts", "A3"))
		assert.Equal(t, "VMWARE", sp.GetCellValue("Hosts", "B3"))
		assert.Equal(t, "Puzzait", sp.GetCellValue("Hosts", "C3"))
		assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", sp.GetCellValue("Hosts", "D3"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU E5630  @ 2.53GHz", sp.GetCellValue("Hosts", "E3"))
		assert.Equal(t, "2", sp.GetCellValue("Hosts", "F3"))
		assert.Equal(t, "1", sp.GetCellValue("Hosts", "G3"))
		assert.Equal(t, "2", sp.GetCellValue("Hosts", "H3"))
		assert.Equal(t, "latest", sp.GetCellValue("Hosts", "I3"))
		assert.Equal(t, "7/1/20 09:18", sp.GetCellValue("Hosts", "J3"))
		assert.Equal(t, "TST", sp.GetCellValue("Hosts", "K3"))
		assert.Equal(t, "Italy", sp.GetCellValue("Hosts", "L2"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "M3"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "N3"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "O3"))
		assert.Equal(t, "Red Hat Enterprise Linux - 7.6", sp.GetCellValue("Hosts", "P3"))
		assert.Equal(t, "0", sp.GetCellValue("Hosts", "Q3"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "R3"))
		assert.Equal(t, "3", sp.GetCellValue("Hosts", "S3"))
		assert.Equal(t, "1", sp.GetCellValue("Hosts", "T3"))
		assert.Equal(t, "", sp.GetCellValue("Hosts", "U3"))
		assert.Equal(t, "0", sp.GetCellValue("Hosts", "V3"))
	})

	t.Run("Db error", func(t *testing.T) {
		db.EXPECT().
			GetHostDataSummaries(filters).
			Return(nil, aerrMock)

		actual, err := as.SearchHostsAsXLSX(filters)
		assert.Nil(t, actual)
		assert.EqualError(t, err, aerrMock.Error())
	})

}

func TestGetHostDataSummaries(t *testing.T) {
	testCases := []struct {
		filters dto.SearchHostsFilters
		res     []dto.HostDataSummary
		err     error
	}{
		{
			filters: dto.SearchHostsFilters{
				Search:      []string{"foo", "bar", "foobarx"},
				SortBy:      "Memory",
				SortDesc:    true,
				Location:    "Italy",
				Environment: "PROD",
				OlderThan:   utils.P("2019-12-05T14:02:03Z"),
				PageNumber:  1,
				PageSize:    1,
			},
			res: []dto.HostDataSummary{
				{
					CreatedAt:               time.Now(),
					Hostname:                "pluto",
					Location:                "Germany",
					Environment:             "TEST",
					AgentVersion:            "0.0.1-alpha",
					Info:                    model.Host{},
					ClusterMembershipStatus: model.ClusterMembershipStatus{},
					Databases:               map[string][]string{},
				},
			},
			err: nil,
		},
		{
			filters: dto.SearchHostsFilters{},
			res:     nil,
			err:     aerrMock,
		},
	}

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		db := NewMockMongoDatabaseInterface(mockCtrl)
		as := APIService{
			Database: db,
		}

		db.EXPECT().GetHostDataSummaries(tc.filters).Return(tc.res, tc.err).Times(1)

		res, err := as.GetHostDataSummaries(tc.filters)
		if tc.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tc.err.Error())
		}

		assert.Equal(t, tc.res, res)
	}
}

func TestGetHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := dto.HostData{
		Alerts: []model.Alert{
			{
				AlertAffectedTechnology: nil,
				AlertCategory:           model.AlertCategoryEngine,
				AlertCode:               "NEW_SERVER",
				AlertSeverity:           "INFO",
				AlertStatus:             "NEW",
				Date:                    utils.P("2020-04-07T08:52:59.871Z"),
				Description:             "The server 'test-virt' was added to ercole",
				OtherInfo: map[string]interface{}{
					"hostname": "test-virt",
				},
				ID: utils.Str2oid("5e8c234b24f648a08585bd42"),
			},
		},
		Archived:    false,
		Cluster:     "Puzzait",
		CreatedAt:   utils.P("2020-04-07T08:52:59.869Z"),
		Environment: "PROD",
		Filesystems: []model.Filesystem{
			{
				AvailableSpace: 4.60000000e+09,
				Filesystem:     "/dev/mapper/vg_os-lv_root",
				MountedOn:      "/",
				Size:           8.00000000e+09,
				Type:           "xfs",
				UsedSpace:      3.50000000e+09,
			},
		},
		History: []model.History{
			{
				CreatedAt: utils.P("2020-04-07T08:52:59.869Z"),
				ID:        utils.Str2oid("5e8c234b24f648a08585bd41"),
			},
		},
		SchemaVersion: 3,
		Hostname:      "test-virt",
		Info: model.Host{
			CPUCores:                      1,
			CPUFrequency:                  "2.50GHz",
			CPUModel:                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			CPUSockets:                    2,
			CPUThreads:                    2,
			CoresPerSocket:                1,
			HardwareAbstraction:           "VIRT",
			HardwareAbstractionTechnology: "VMWARE",
			Hostname:                      "test-virt",
			Kernel:                        "Linux",
			KernelVersion:                 "3.10.0-862.9.1.el7.x86_64",
			MemoryTotal:                   3,
			OS:                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			OSVersion:                     "7.5",
			SwapTotal:                     4,
			ThreadsPerCore:                2,
		},
		Features: model.Features{
			Oracle: &model.OracleFeature{
				Database: &model.OracleDatabaseFeature{
					Databases: []model.OracleDatabase{
						{
							Licenses: []model.OracleDatabaseLicense{
								{
									LicenseTypeID: "A90611",
									Count:         2,
								},
								{
									LicenseTypeID: "L76084",
									Name:          "Real Application Clusters One Node",
									Count:         2,
								},
								{
									LicenseTypeID: "A90619",
									Name:          "Real Application Clusters",
									Count:         1,
								},
							},
						},
					},
				},
			},
		},
		Location:            "Italy",
		VirtualizationNode:  "s157-cb32c10a56c256746c337e21b3f82402",
		ServerSchemaVersion: 1,
		ServerVersion:       "latest",
		AgentVersion:        "1.6.1",
		ID:                  utils.Str2oid("5e8c234b24f648a08585bd41"),
	}

	db.EXPECT().GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	).Return(&expectedRes, nil).Times(1)

	res, err := as.GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	)
	require.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
}

func TestGetHost_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	)
	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestListLocations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []string{
		"Italy",
		"Germany",
	}

	db.EXPECT().ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListAllLocations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []string{
		"Italy",
		"Germany",
	}

	db.EXPECT().ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListAllLocations_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.ListAllLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestListEnvironments_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []string{
		"TST",
		"SVIL",
		"PROD",
	}

	db.EXPECT().ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListEnvironments_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestDismissHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	alwaysTheSameMoment := func() time.Time {
		return time.Date(1994, 11, 19, 0, 3, 3, 0, time.Local)
	}
	as := APIService{
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        alwaysTheSameMoment,
	}

	var count int64

	expectedRes := []map[string]interface{}{
		{
			"hostname": "foobar",
		},
	}

	listContracts := []dto.OracleDatabaseContractFE{}

	filter := dto.AlertsFilter{OtherInfo: map[string]interface{}{"hostname": "foobar"}}
	db.EXPECT().RemoveAlertsNODATA(filter).Return(nil).Times(1)
	db.EXPECT().CountAlertsNODATA(filter).Return(count, nil).Times(1)
	db.EXPECT().UpdateAlertsStatus(filter, model.AlertStatusAck).Return(nil)
	db.EXPECT().UpdateAlertsStatus(filter, model.AlertStatusDismissed).Return(nil)
	commonFilters := dto.NewSearchHostsFilters()
	db.EXPECT().SearchHosts(
		"hostnames",
		commonFilters,
	).Return(expectedRes, nil)
	db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).Return(listContracts, nil)

	asc.EXPECT().ThrowNewAlert(gomock.Any()).Return(nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.TechnologyOracleDatabasePtr, alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCategoryEngine, alert.AlertCategory)
		assert.Equal(t, model.AlertCodeDismissHost, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityInfo, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, as.TimeNow(), alert.Date)
		assert.Equal(t, "Host foobar was dismissed", alert.Description)
	})

	db.EXPECT().DismissHost("foobar").Return(nil).Times(1)

	err := as.DismissHost("foobar")
	require.NoError(t, err)
}

func TestDismissHost_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	var count int64

	expectedRes := []map[string]interface{}{
		{
			"hostname": "foobar",
		},
	}

	listContracts := []dto.OracleDatabaseContractFE{}

	filter := dto.AlertsFilter{OtherInfo: map[string]interface{}{"hostname": "foobar"}}
	db.EXPECT().RemoveAlertsNODATA(filter).Return(nil).Times(1)
	db.EXPECT().CountAlertsNODATA(filter).Return(count, nil).Times(1)
	db.EXPECT().UpdateAlertsStatus(filter, model.AlertStatusAck).Return(nil)
	db.EXPECT().UpdateAlertsStatus(filter, model.AlertStatusDismissed).Return(nil)
	commonFilters := dto.NewSearchHostsFilters()
	db.EXPECT().SearchHosts(
		"hostnames",
		commonFilters,
	).Return(expectedRes, nil)
	db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).Return(listContracts, nil)
	db.EXPECT().DismissHost("foobar").Return(aerrMock).Times(1)

	err := as.DismissHost("foobar")
	assert.Error(t, err)
}

func TestGetMissingDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	as := APIService{
		Database:       db,
		AlertSvcClient: asc,
	}

	expected := []dto.OracleDatabaseMissingDbs{
		{
			Hostname: "host01",
			MissingDatabases: []model.MissingDatabase{
				{
					Name: "db01",
					Ignorable: model.Ignorable{
						Ignored:        false,
						IgnoredComment: "",
					},
				},
				{
					Name: "db02",
					Ignorable: model.Ignorable{
						Ignored:        true,
						IgnoredComment: "this is no longer needed",
					},
				},
			},
		},
	}

	db.EXPECT().GetMissingDatabases().Return(expected, nil)

	actual, err := as.GetMissingDatabases()

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestGetVirtualHostWithoutCluster_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	as := APIService{
		Database:       db,
		AlertSvcClient: asc,
	}

	expected := []dto.VirtualHostWithoutCluster{
		{
			Hostname:                      "host01",
			HardwareAbstractionTechnology: "VIRT",
		},
	}

	db.EXPECT().FindVirtualHostWithoutCluster().Return(expected, nil)

	actual, err := as.GetVirtualHostWithoutCluster()

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestListLocationsLicenses(t *testing.T) {
	testCases := []struct {
		name          string
		config        config.Configuration
		user          model.User
		buildStubs    func(db *MockMongoDatabaseInterface)
		checkResponse func(t *testing.T, res []string, err error)
	}{
		{
			name: "Ok user",
			config: config.Configuration{
				APIService: config.APIService{
					ScopeAsLocation: "loc1,loc2,loc3",
				},
			},
			user: model.User{
				Username: "user01",
			},
			buildStubs: func(db *MockMongoDatabaseInterface) {
				db.EXPECT().GetUserLocations("user01").
					Return([]string{"loc1", "loc2", "loc3", "loc4"}, nil)
			},

			checkResponse: func(t *testing.T, res []string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedResult := []string{"loc1,loc2,loc3", "loc4"}
				assert.Equal(t, expectedResult, res)
			},
		},
		{
			name: "Ok Admin",
			config: config.Configuration{
				APIService: config.APIService{
					ScopeAsLocation: "loc1,loc2,loc3",
				},
			},
			user: model.User{
				Username: "user01",
				Groups:   []string{"admin"},
			},
			buildStubs: func(db *MockMongoDatabaseInterface) {
				db.EXPECT().ListAllLocations("", "", utils.MAX_TIME).
					Return([]string{"loc1", "loc2", "loc3", "loc4"}, nil)
			},

			checkResponse: func(t *testing.T, res []string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedResult := []string{"loc1,loc2,loc3", "loc4"}
				assert.Equal(t, expectedResult, res)
			},
		},
		{
			name: "ScopeAsLocation empty",
			config: config.Configuration{
				APIService: config.APIService{
					ScopeAsLocation: "",
				},
			},
			user: model.User{
				Username: "user01",
			},
			buildStubs: func(db *MockMongoDatabaseInterface) {
				db.EXPECT().GetUserLocations("user01").
					Return([]string{"loc1", "loc2", "loc3", "loc4"}, nil)
			},

			checkResponse: func(t *testing.T, res []string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				expectedResult := []string{"loc1", "loc2", "loc3", "loc4"}
				assert.Equal(t, expectedResult, res)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			db := NewMockMongoDatabaseInterface(mockCtrl)

			as := APIService{
				Database: db,
				Config:   tc.config,
			}

			tc.buildStubs(db)

			res, err := as.ListLocationsLicenses(tc.user)

			tc.checkResponse(t, res, err)
		})
	}
}
