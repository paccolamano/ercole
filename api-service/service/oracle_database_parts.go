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

package service

import (
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (as *APIService) GetOracleDatabaseAgreementPartsList() ([]model.OracleDatabasePart, utils.AdvancedErrorInterface) {
	parts, err := as.Database.GetOracleDatabaseParts()
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return parts, nil
}

// GetOracleDatabaseAgreementPartsList return Oracle/Database agreement parts mapped by their PartID
func (as *APIService) GetOracleDatabaseAgreementPartsMap() (map[string]model.OracleDatabasePart, utils.AdvancedErrorInterface) {
	parts, err := as.Database.GetOracleDatabaseParts() // Should call GetOracleDatabaseAgreementPartsList ?
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	partsMap := make(map[string]model.OracleDatabasePart)
	for _, part := range parts {
		partsMap[part.PartID] = part
	}

	return partsMap, nil
}

// GetOraclePart return a Part by ID
func (as *APIService) GetOraclePart(partID string) (*model.OracleDatabasePart, utils.AdvancedErrorInterface) {
	parts, err := as.Database.GetOracleDatabaseParts()
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	for _, part := range parts {
		if partID == part.PartID {
			return &part, nil
		}
	}

	return nil, utils.AerrOracleDatabaseAgreementInvalidPartID
}

func (as *APIService) GetOracleDatabaseLicensesCompliance() ([]dto.OracleDatabaseLicenseUsage, utils.AdvancedErrorInterface) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	err2 := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	if err2 != nil {
		return nil, utils.NewAdvancedErrorPtr(err2, "DB ERROR")
	}

	licenses := make(map[string]*dto.OracleDatabaseLicenseUsage)

	for _, host := range hosts {
		license, ok := licenses[host.PartID]
		if !ok {
			license = &dto.OracleDatabaseLicenseUsage{
				PartID: host.PartID,
			}

			licenses[license.PartID] = license
		}

		license.Consumed += host.OriginalCount
	}

	for _, agreement := range agreements {
		license, ok := licenses[agreement.PartID]
		if !ok {
			license = &dto.OracleDatabaseLicenseUsage{
				PartID: agreement.PartID,
			}

			licenses[license.PartID] = license
		}

		if agreement.Unlimited {
			license.Unlimited = true
		}

		for _, host := range agreement.Hosts {
			license.Covered += host.CoveredLicensesCount
		}
	}

	parts, err := as.GetOracleDatabaseAgreementPartsMap()
	if err != nil {
		return nil, err
	}

	result := make([]dto.OracleDatabaseLicenseUsage, 0, len(licenses))
	for _, license := range licenses {
		if license.Consumed == 0 {
			license.Compliance = 1
		} else {
			license.Compliance = license.Covered / license.Consumed
		}

		license.ItemDescription = parts[license.PartID].ItemDescription
		license.Metric = parts[license.PartID].Metric

		result = append(result, *license)
	}

	return result, nil
}
