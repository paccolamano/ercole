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
	"errors"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListOracleDiskGroups_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleDatabaseDiskGroupDto{
		{
			Hostname: "hostname",
			Databases: []string{
				"db1",
				"db2",
			},
			OracleDatabaseDiskGroup: model.OracleDatabaseDiskGroup{
				DiskGroupName: "disk one",
				TotalSpace:    100,
				FreeSpace:     75,
			},
		},
	}
	db.EXPECT().ListOracleDiskGroups(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expected, nil)

	res, err := as.ListOracleDiskGroups(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestListOracleDiskGroups_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	e := errors.New("connection error")
	db.EXPECT().ListOracleDiskGroups(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(nil, e)

	res, err := as.ListOracleDiskGroups(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
	require.ErrorIs(t, err, e)
	assert.Nil(t, res)
}
