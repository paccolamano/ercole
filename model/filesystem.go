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

package model

import (
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
)

// Filesystem holds information about mounted filesystem and used space
type Filesystem struct {
	Filesystem     string                 `bson:"Filesystem"`
	Type           string                 `bson:"Type"`
	Size           int64                  `bson:"Size"`
	UsedSpace      int64                  `bson:"UsedSpace"`
	AvailableSpace int64                  `bson:"AvailableSpace"`
	MountedOn      string                 `bson:"MountedOn"`
	OtherInfo      map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Filesystem) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Filesystem) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// FilesystemBsonValidatorRules contains mongodb validation rules for Filesystem
var FilesystemBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Filesystem",
		"Type",
		"Size",
		"UsedSpace",
		"AvailableSpace",
		"MountedOn",
	},
	"properties": bson.M{
		"Filesystem": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"Type": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"Size": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"UsedSpace": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"AvailableSpace": bson.M{
			"bsonType": "number",
		},
		"MountedOn": bson.M{
			"bsonType": "string",
		},
	},
}
