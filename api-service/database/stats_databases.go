// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
func (md *MongoDatabase) GetDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroup(bson.M{
				"_id":   "$environment",
				"count": mu.APOSum(mu.APOCond("$extra.databases", mu.APOSize("$extra.databases"), 0)),
			}),
			mu.APProject(bson.M{
				"_id":         false,
				"environment": "$_id",
				"count":       true,
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabaseVersionStats return a array containing the number of databases per version
func (md *MongoDatabase) GetDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroupAndCountStages("version", "count", "$database.version"),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (md *MongoDatabase) GetTopReclaimableDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APProject(bson.M{
				"hostname": true,
				"dbname":   "$database.name",
				"reclaimable_segment_advisors": mu.APOSumReducer("$database.segment_advisors",
					mu.APOConvertErrorable("$$this.reclaimable", "double", 0.5),
				),
			}),
			mu.APSort(bson.M{
				"reclaimable_segment_advisors": -1,
			}),
			mu.APLimit(limit),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetTopWorkloadDatabaseStats return a array containing top databases by workload
func (md *MongoDatabase) GetTopWorkloadDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APProject(bson.M{
				"hostname": true,
				"dbname":   "$database.name",
				"workload": mu.APOConvertToDoubleOrZero("$database.work"),
			}),
			mu.APSort(bson.M{
				"workload": -1,
			}),
			mu.APLimit(limit),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabasePatchStatusStats return a array containing the number of databases per patch status
func (md *MongoDatabase) GetDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APProject(bson.M{
				"database.last_psus": mu.APOReduce(
					mu.APOMap("$database.last_psus", "psu", mu.APOMergeObjects(
						"$$psu",
						bson.M{
							"date": mu.APODateFromString("$$psu.date", "%Y-%m-%d"),
						},
					)),
					nil,
					mu.APOCond(
						mu.APOEqual("$$value", nil),
						"$$this",
						mu.APOMaxWithCmpExpr("$$value.date", "$$this.date", "$$value", "$$this"),
					),
				),
			}),
			mu.APGroupAndCountStages("status", "id",
				mu.APOCond(mu.APOGreater("$database.last_psus.date", windowTime), "OK", "KO"),
			),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (md *MongoDatabase) GetDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroupAndCountStages("dataguard", "count", "$database.dataguard"),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (md *MongoDatabase) GetDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroupAndCountStages("rac", "count", mu.APOAny("$database.features", "fe",
				mu.APOAnd(
					mu.APOEqual("$$fe.name", "Real Application Clusters"),
					mu.APOEqual("$$fe.status", true),
				),
			)),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (md *MongoDatabase) GetDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroupAndCountStages("archivelog", "count",
				mu.APOEqual("$database.archive_log", "ARCHIVELOG"),
			),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetTotalDatabaseWorkStats return the total work of databases
func (md *MongoDatabase) GetTotalDatabaseWorkStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float32

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum(mu.APOConvertToDoubleOrZero("$database.work")),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out["value"], nil
}

// GetTotalDatabaseMemorySizeStats return the total of memory size of databases
func (md *MongoDatabase) GetTotalDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id": 0,
				"value": mu.APOSum(mu.APOAdd(
					mu.APOConvertToDoubleOrZero("$database.pga_target"),
					mu.APOConvertToDoubleOrZero("$database.sga_target"),
					mu.APOConvertToDoubleOrZero("$database.memory_target"),
				)),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return float32(out["value"]), nil
}

// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases
func (md *MongoDatabase) GetTotalDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum(mu.APOConvertToDoubleOrZero("$database.used")),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return float32(out["value"]), nil
}

// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases
func (md *MongoDatabase) GetTotalDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum(mu.APOConvertToDoubleOrZero("$database.segments_size")),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return float32(out["value"]), nil
}

// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
func (md *MongoDatabase) GetDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{} = map[string]interface{}{
		"count":     0,
		"used":      0,
		"compliant": true,
	}

	//Find the informations
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APLookupPipeline("hosts", bson.M{
				"license_name": "$_id",
			}, "used", mu.MAPipeline(
				mu.APMatch(bson.M{
					"archived": false,
				}),
				mu.APProject(bson.M{
					"hostname": 1,
					"databases": mu.APOReduce(
						mu.APOFilter(
							mu.APOMap("$extra.databases", "db", bson.M{
								"name": "$$db.name",
								"count": mu.APOLet(
									bson.M{
										"val": mu.APOArrayElemAt(mu.APOFilter("$$db.licenses", "lic", mu.APOEqual("$$lic.name", "$$license_name")), 0),
									},
									"$$val.count",
								),
							}),
							"db",
							mu.APOGreater("$$db.count", 0),
						),
						bson.M{"count": 0, "dbs": bson.A{}},
						bson.M{
							"count": mu.APOMax("$$value.count", "$$this.count"),
							"dbs": bson.M{
								"$concatArrays": bson.A{
									"$$value.dbs",
									bson.A{"$$this.name"},
								},
							},
						},
					),
				}),
				mu.APMatch(bson.M{
					"databases.count": bson.M{
						"$gt": 0,
					},
				}),
				mu.APLookupSimple("currentClusters", "hostname", "cluster.vms.hostname", "cluster"),
				mu.APSet(bson.M{
					"cluster": mu.APOArrayElemAt("$cluster", 0),
				}),
				// mu.APSet(bson.M{
				// 	"cluster": mu.APOArrayElemAt(
				// 		mu.APOFilter("$cluster.cluster.vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")),
				// 		0,
				// 	),
				// }),
				mu.APSet(bson.M{
					"cluster_name": "$cluster.cluster.name",
					"cluster_cpu":  "$cluster.cluster.cpu",
				}),
				mu.APUnset("cluster"),
				mu.APGroup(bson.M{
					"_id": mu.APOCond(
						"$cluster_name",
						mu.APOConcat("cluster_§$#$§_", "$cluster_name"),
						mu.APOConcat("hostname_§$#$§_", "$hostname"),
					),
					"license":     mu.APOMaxAggr("$databases.count"),
					"cluster_cpu": mu.APOMaxAggr("$cluster_cpu"),
				}),
				mu.APSet(bson.M{
					"license": mu.APOCond(
						"$cluster_cpu",
						mu.APODivide("$cluster_cpu", 2),
						"$license",
					),
				}),
				mu.APGroup(bson.M{
					"_id":   0,
					"value": mu.APOSum("$license"),
				}),
			)),
			mu.APSet(bson.M{
				"used": mu.APOArrayElemAt("$used", 0),
			}),
			mu.APSet(bson.M{
				"used": mu.APOIfNull(mu.APOCeil("$used.value"), 0),
			}),
			mu.APSet(bson.M{
				"compliance": mu.APOGreaterOrEqual("$count", "$used"),
			}),
			mu.APGroup(bson.M{
				"_id":                       0,
				"licenses_number":           mu.APOSum(1),
				"count":                     mu.APOSum("$count"),
				"used":                      mu.APOSum("$used"),
				"compliant_licenses_number": mu.APOSum(mu.APOCond("$compliance", 1, 0)),
			}),
			mu.APProject(bson.M{
				"_id":       0,
				"count":     1,
				"used":      1,
				"compliant": mu.APOEqual("$licenses_number", "$compliant_licenses_number"),
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return out, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
