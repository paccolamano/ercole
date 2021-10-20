// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migrations

import (
	"context"

	"github.com/ercole-io/ercole/v2/utils"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		if err := migrateFromCatchAllToBasket(db); err != nil {
			return err
		}
		if err := migrateOracleLicensesAgreement(db); err != nil {
			return err
		}

		return nil

	}, func(db *mongo.Database) error {
		return nil
	})
}

func migrateFromCatchAllToBasket(db *mongo.Database) error {
	collection := "oracle_database_agreements"
	ctx := context.TODO()

	filter := bson.M{}
	update := bson.M{"$rename": bson.M{"catchAll": "basket"}}

	_, err := db.Collection(collection).UpdateMany(ctx, filter, update)
	if err != nil {
		return utils.NewError(err, "Can't rename from catchAll to basket", collection)
	}

	return nil
}

func migrateOracleLicensesAgreement(db *mongo.Database) error {
	collection := "oracle_database_agreements"
	ctx := context.TODO()

	filter := bson.M{"unlimited": true}
	update := bson.M{"$set": bson.M{"basket": true}}

	_, err := db.Collection(collection).UpdateMany(ctx, filter, update)
	if err != nil {
		return utils.NewError(err, "Can't update basket to true with unlimited = true", collection)
	}

	return nil
}