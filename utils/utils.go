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

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MIN_TIME time.Time = time.Unix(0, 0)
var MAX_TIME time.Time = time.Now().AddDate(1000, 0, 0)

// ToJSON convert v to a string containing the equivalent json rappresentation
func ToJSON(v interface{}) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

// ToMongoJSON convert v to a string containing the equivalent json rappresentation
func ToMongoJSON(v interface{}) string {
	raw, err := bson.MarshalExtJSON(v, false, false)
	if err != nil {
		panic(err)
	}
	return string(raw)
}

// FromJSON convert a json str to interface containing the equivalent json rappresentation
func FromJSON(str []byte) interface{} {
	var out map[string]interface{}
	json.Unmarshal(str, &out)
	return out
}

// ToJSONMongoCursor extract all items from a cursors and return its json rappresentation
func ToJSONMongoCursor(cur *mongo.Cursor) string {
	var out = make([]map[string]interface{}, 0)
	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if err := cur.Decode(&item); err != nil {
			panic(err)
		}
		out = append(out, item)
	}
	return ToIdentedJSON(out)
}

// ToIdentedJSON convert v to a string containing the equivalent idented json rappresentation
func ToIdentedJSON(v interface{}) string {
	raw, _ := json.MarshalIndent(v, "", "  ")
	return string(raw)
}

// Intptr return a point to the int passed in the argument
func Intptr(v int64) *int64 {
	return &v
}

// Contains return true if a contains x, otherwise false.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// ContainsI return true if a contains x (insensitive-case), otherwise false.
func ContainsI(a []string, x string) bool {
	for _, n := range a {
		if strings.EqualFold(x, n) {
			return true
		}
	}
	return false
}

// Difference returns the elements in `a` that aren't in `b`
// If a has multiple times on item, which is in b even only once, no occurrences will be returned
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}

	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

// Remove return slice without element at position i, mantaining order
func Remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

func RemoveString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Str2bool parse a string to a boolean
func Str2bool(in string, defaultValue bool) (bool, error) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseBool(in); err != nil {
		return false, NewError(err, "Unable to parse string to bool")
	} else {
		return val, nil
	}
}

// Str2int parse a string to a int
func Str2int(in string, defaultValue int) (int, error) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseInt(in, 10, 32); err != nil {
		return -1, NewError(err, "Unable to parse string to int")
	} else {
		return int(val), nil
	}
}

// Str2float64 parse a string to a float64
func Str2float64(in string, defaultValue float64) (float64, error) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseFloat(in, 32); err != nil {
		return -1, NewError(err, "Unable to parse string to float")
	} else {
		return float64(val), nil
	}
}

// Str2time parse a string to a time
func Str2time(in string, defaultValue time.Time) (time.Time, error) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := time.Parse(time.RFC3339, in); err != nil {
		return time.Time{}, NewError(err, "Unable to parse string to time.Time")
	} else {
		return val, nil
	}
}

// Str2ptr return a pointer to a copy of s
// Go pass values by copy, so a copy of the value passed
// is already made when calling method
func Str2ptr(s string) *string {
	return &s
}

// NewAPIUrl return a new url crafted using the parameters
func NewAPIUrl(baseURL string, username string, password string, path string, params url.Values) *url.URL {
	u := NewAPIUrlNoParams(baseURL, username, password, path)
	u.RawQuery = params.Encode()

	return u
}

// NewAPIUrlNoParams return a new url crafted using the parameters
func NewAPIUrlNoParams(baseURL string, username string, password string, path string) *url.URL {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	u.User = url.UserPassword(username, password)
	u.Path += path

	return u
}

// FindNamedMatches return the map of the groups of str
func FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

// DownloadFile download the file from url into the filepath
func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	client := http.Client{Timeout: 1 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// IsVersionLessThan return true if a is a version less than b
func IsVersionLessThan(a, b string) (bool, error) {
	va, err := version.NewVersion(a)
	if err != nil {
		return false, err
	}

	vb, err := version.NewVersion(b)
	if err != nil {
		return false, err
	}

	return va.LessThan(vb), nil
}

func IsVersionEqual(a, b string) (bool, error) {
	thisSemver, err := version.NewSemver(a)
	if err != nil {
		return false, err
	}

	semver, err := version.NewSemver(b)
	if err != nil {
		return false, err
	}

	return thisSemver.Equal(semver), nil
}

func HideMongoDBPassword(uri string) string {
	m := regexp.MustCompile(`\/\/([^:]+):(.*)@`)
	return m.ReplaceAllString(uri, "//***:***@")
}

// TruncateFloat64 truncate float in 2 decimals
func TruncateFloat64(src float64) float64 {
	return float64(int(src*100)) / 100
}

func RemoveDuplicate[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ContainsSomeI(slice []string, items ...string) bool {
	for _, n := range slice {
		if ContainsI(items, n) {
			return true
		}
	}
	return false
}
