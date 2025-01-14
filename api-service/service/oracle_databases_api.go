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

// Package service is a package that provides methods for querying data
package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/utils/exutils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ercole-io/ercole/v2/api-service/dto"
)

// SearchOracleDatabaseAddms search addms
func (as *APIService) SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error) {
	return as.Database.SearchOracleDatabaseAddms(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// SearchOracleDatabaseSegmentAdvisors search segment advisors
func (as *APIService) SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]dto.OracleDatabaseSegmentAdvisor, error) {
	result := make([]dto.OracleDatabaseSegmentAdvisor, 0)

	if dbSegmentAdvisors, err := as.Database.SearchOracleDatabaseSegmentAdvisors(strings.Split(search, " "), sortBy, sortDesc, location, environment, olderThan); err != nil {
		return nil, err
	} else {
		result = append(result, dbSegmentAdvisors...)
	}

	if pdbSegmentAdvisors, err := as.Database.SearchOraclePdbSegmentAdvisors(sortBy, sortDesc, location, environment, olderThan); err != nil {
		return nil, err
	} else {
		result = append(result, pdbSegmentAdvisors...)
	}

	return result, nil
}

func (as *APIService) SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	segmentAdvisors, err := as.SearchOracleDatabaseSegmentAdvisors("", "", false, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	sheet := "Segment_Advisor"
	headers := []string{
		"ReclaimableGB",
		"GB Total",
		"Retrieve",
		"Hostname",
		"DB Names",
		"PDB Name",
		"Segment Owner",
		"Segment Name",
		"Segment Type",
		"Partition Name",
		"Recommendation",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range segmentAdvisors {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Reclaimable)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentsSize)

		if val.SegmentsSize == 0 {
			nextAxis()
		} else {
			sheets.SetCellValue(sheet, nextAxis(), val.Reclaimable/val.SegmentsSize)
		}

		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.Dbname)
		sheets.SetCellValue(sheet, nextAxis(), val.PdbName)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentOwner)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentName)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentType)
		sheets.SetCellValue(sheet, nextAxis(), val.PartitionName)
		sheets.SetCellValue(sheet, nextAxis(), val.Recommendation)
	}

	return sheets, err
}

// SearchOracleDatabasePatchAdvisors search patch advisors
func (as *APIService) SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) (*dto.PatchAdvisorResponse, error) {
	return as.Database.SearchOracleDatabasePatchAdvisors(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, windowTime, location, environment, olderThan, status)
}

func (as *APIService) SearchOracleDatabasePatchAdvisorsAsXLSX(windowTime time.Time, filter dto.GlobalFilter) (*excelize.File, error) {
	patchAdvisorResponse, err := as.Database.SearchOracleDatabasePatchAdvisors([]string{}, "", false, -1, -1, windowTime, filter.Location, filter.Environment, filter.OlderThan, "")
	if err != nil {
		return nil, err
	}

	sheet := "Patch_Advisor"
	headers := []string{
		"Hostname",
		"Database",
		"Version",
		"Release Date",
		"PSU",
		"Status",
		"4 Months",
		"6 Months",
		"12 Months",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range patchAdvisorResponse.Content {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.Hostname)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.DbName)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.Dbver)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.Date.Time().UTC().String())
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.Description)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.Status)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.FourMonths)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.SixMonths)
		sheets.SetCellValue("Patch_Advisor", nextAxis(), val.TwelveMonths)
	}

	return sheets, err
}

// SearchOracleDatabases search databases
func (as *APIService) SearchOracleDatabases(f dto.SearchOracleDatabasesFilter) (*dto.OracleDatabaseResponse, error) {
	response, err := as.Database.SearchOracleDatabases(strings.Split(f.Search, " "), f.SortBy, f.SortDesc,
		f.PageNumber, f.PageSize, f.Location, f.Environment, f.OlderThan)
	if err != nil {
		return nil, err
	}

	// for i := range response.Content {
	// 	color, err := as.GetOraclePsqlMigrabilitiesSemaphore(response.Content[i].Hostname, response.Content[i].Name)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	response.Content[i].PgsqlMigrabilitySemaphore = color

	// 	canbemigrate, err := as.CanMigrateLicense(response.Content[i].Hostname, response.Content[i].Name, f.GlobalFilter)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	response.Content[i].Canbemigrate = canbemigrate
	// }

	return response, nil
}

func (as *APIService) SearchOracleDatabasesAsXLSX(filter dto.SearchOracleDatabasesFilter) (*excelize.File, error) {
	databases, err := as.Database.SearchOracleDatabases(strings.Split(filter.Search, " "),
		filter.SortBy, filter.SortDesc,
		-1, -1,
		filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	convertedColor := map[string]string{
		"red":    "1%",
		"yellow": "50%",
		"green":  "100%",
	}

	for i := range databases.Content {
		color, err := as.GetOraclePsqlMigrabilitiesSemaphore(databases.Content[i].Hostname, databases.Content[i].Name)
		if err != nil {
			return nil, err
		}

		databases.Content[i].PgsqlMigrabilitySemaphore = convertedColor[color]
	}

	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_databases.xlsx")
	if err != nil {
		return nil, err
	}

	for i, val := range databases.Content {
		i += 2 // offset for headers
		file.SetCellValue("Databases", fmt.Sprintf("A%d", i), val.Name)
		file.SetCellValue("Databases", fmt.Sprintf("B%d", i), val.UniqueName)
		file.SetCellValue("Databases", fmt.Sprintf("C%d", i), val.DbID)
		file.SetCellValue("Databases", fmt.Sprintf("D%d", i), val.Hostname)
		file.SetCellValue("Databases", fmt.Sprintf("E%d", i), val.Version)
		file.SetCellValue("Databases", fmt.Sprintf("F%d", i), val.Status)
		file.SetCellValue("Databases", fmt.Sprintf("G%d", i), val.Environment)
		file.SetCellValue("Databases", fmt.Sprintf("H%d", i), val.Location)
		file.SetCellValue("Databases", fmt.Sprintf("I%d", i), val.Charset)
		file.SetCellValue("Databases", fmt.Sprintf("J%d", i), val.BlockSize)
		file.SetCellValue("Databases", fmt.Sprintf("K%d", i), val.CPUCount)

		if val.Work != nil {
			file.SetCellValue("Databases", fmt.Sprintf("L%d", i), *val.Work)
		} else {
			file.SetCellValue("Databases", fmt.Sprintf("L%d", i), "")
		}

		file.SetCellValue("Databases", fmt.Sprintf("M%d", i), val.Memory)
		file.SetCellValue("Databases", fmt.Sprintf("N%d", i), val.DatafileSize)
		file.SetCellValue("Databases", fmt.Sprintf("O%d", i), val.SegmentsSize)
		file.SetCellValue("Databases", fmt.Sprintf("P%d", i), val.Archivelog)
		file.SetCellValue("Databases", fmt.Sprintf("Q%d", i), val.Dataguard)
		file.SetCellValue("Databases", fmt.Sprintf("R%d", i), val.Rac)
		file.SetCellValue("Databases", fmt.Sprintf("S%d", i), val.Ha)

		if val.IsCDB {
			file.SetCellValue("Databases", fmt.Sprintf("T%d", i), "True")

			pdbnames := make([]string, 0)

			for _, p := range val.PDBs {
				pdbnames = append(pdbnames, p.Name)
			}

			file.SetCellValue("Databases", fmt.Sprintf("U%d", i), strings.Join(pdbnames, ","))
		} else {
			file.SetCellValue("Databases", fmt.Sprintf("T%d", i), "False")
			file.SetCellValue("Databases", fmt.Sprintf("U%d", i), "")
		}

		file.SetCellValue("Databases", fmt.Sprintf("V%d", i), val.PgsqlMigrabilitySemaphore)
		file.SetCellValue("Databases", fmt.Sprintf("W%d", i), val.ClusterwareVersion)
		file.SetCellValue("Databases", fmt.Sprintf("X%d", i), val.Role)
		file.SetCellValue("Databases", fmt.Sprintf("Y%d", i), strings.Join(val.PoliciesAudit, ","))

		serviceNames := make([]string, 0)

		for _, s := range val.Services {
			if s.Name != nil {
				serviceNames = append(serviceNames, *s.Name)
			}
		}

		file.SetCellValue("Databases", fmt.Sprintf("Z%d", i), strings.Join(serviceNames, ","))
	}

	return file, nil
}

// SearchOracleDatabaseUsedLicenses return the list of used licenses
func (as *APIService) SearchOracleDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.OracleDatabaseUsedLicenseSearchResponse, error) {
	return as.Database.SearchOracleDatabaseUsedLicenses(hostname, sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}
