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

package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/api-service/auth"
)

// GetApiControllerHandler setup the routes of the router using the handler in the controller as http handler
func (ctrl *APIController) GetApiControllerHandler(auth auth.AuthenticationProvider) http.Handler {
	router := mux.NewRouter()

	//Add the routes
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Pong")); err != nil {
			ctrl.Log.Error(err)
			return
		}
	})

	router.HandleFunc("/user/login", auth.GetToken).Methods("POST")

	//Enable authentication using the ctrl
	subrouter := router.NewRoute().Subrouter()
	subrouter.Use(auth.AuthenticateMiddleware)
	ctrl.setupProtectedRoutes(subrouter)

	return router
}

func (ctrl *APIController) setupProtectedRoutes(router *mux.Router) {
	// ERCOLE
	router.HandleFunc("/version", ctrl.GetVersion).Methods("GET")

	// HOSTS
	router.HandleFunc("/hosts", ctrl.SearchHosts).Methods("GET")
	router.HandleFunc("/hosts/count", ctrl.GetHostsCountStats).Methods("GET")
	router.HandleFunc("/hosts/environments/frequency", ctrl.GetEnvironmentStats).Methods("GET")
	router.HandleFunc("/hosts/types", ctrl.GetTypeStats).Methods("GET")
	router.HandleFunc("/hosts/operating-systems", ctrl.GetOperatingSystemStats).Methods("GET")
	router.HandleFunc("/hosts/locations", ctrl.ListLocations).Methods("GET")
	router.HandleFunc("/hosts/environments", ctrl.ListEnvironments).Methods("GET")
	router.HandleFunc("/hosts/clusters", ctrl.SearchClusters).Methods("GET")
	router.HandleFunc("/hosts/clusters/{name}", ctrl.GetCluster).Methods("GET")

	router.HandleFunc("/hosts/{hostname}", ctrl.GetHost).Methods("GET")
	router.HandleFunc("/hosts/{hostname}", ctrl.DismissHost).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/technologies/oracle/databases/{dbname}/licenses/{licenseTypeID}/ignored/{ignored}", ctrl.UpdateLicenseIgnoredField).Methods("PUT")

	router.HandleFunc("/hosts/technologies", ctrl.ListTechnologies).Methods("GET")

	// ALL TECHNOLOGIES
	router.HandleFunc("/hosts/technologies/all/databases", ctrl.SearchDatabases).Methods("GET")
	router.HandleFunc("/hosts/technologies/all/databases/statistics", ctrl.GetDatabasesStatistics).Methods("GET")
	router.HandleFunc("/hosts/technologies/all/databases/licenses-used", ctrl.GetUsedLicensesPerDatabases).Methods("GET")
	router.HandleFunc("/hosts/{hostname}/technologies/all/databases/licenses-used", ctrl.GetUsedLicensesPerDatabasesByHost).Methods("GET")
	router.HandleFunc("/hosts/technologies/all/databases/licenses-used-per-host", ctrl.GetUsedLicensesPerHost).Methods("GET")
	router.HandleFunc("/hosts/technologies/all/databases/licenses-used-per-cluster", ctrl.GetUsedLicensesPerCluster).Methods("GET")
	router.HandleFunc("/hosts/technologies/all/databases/licenses-compliance", ctrl.GetDatabaseLicensesCompliance).Methods("GET")

	router.HandleFunc("/hosts/technologies/all/databases/grant-dba", ctrl.ListOracleGrantDbaByHostname).Methods("GET")

	// ORACLE
	router.HandleFunc("/hosts/technologies/oracle/databases", ctrl.SearchOracleDatabases).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-unused-instance-resource", ctrl.GetTopUnusedOracleDatabaseInstanceResourceStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/environments", ctrl.GetOracleDatabaseEnvironmentStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/high-reliability", ctrl.GetOracleDatabaseHighReliabilityStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/versions", ctrl.GetOracleDatabaseVersionStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-reclaimable", ctrl.GetTopReclaimableOracleDatabaseStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/patch-status", ctrl.GetOracleDatabasePatchStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-workload", ctrl.GetTopWorkloadOracleDatabaseStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/dataguard-status", ctrl.GetOracleDatabaseDataguardStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/archivelog-status", ctrl.GetOracleDatabaseArchivelogStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/rac-status", ctrl.GetOracleDatabaseRACStatusStats).Methods("GET")

	router.HandleFunc("/hosts/technologies/oracle/databases/statistics", ctrl.GetOracleDatabasesStatistics).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/consumed-licenses", ctrl.SearchOracleDatabaseUsedLicenses).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/licenses-compliance", ctrl.GetOracleDatabaseLicensesCompliance).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/addms", ctrl.SearchOracleDatabaseAddms).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/segment-advisors", ctrl.SearchOracleDatabaseSegmentAdvisors).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/patch-advisors", ctrl.SearchOracleDatabasePatchAdvisors).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/patch-list", ctrl.GetOraclePatchList).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/option-list", ctrl.GetOracleOptionList).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/tablespaces", ctrl.ListOracleDatabaseTablespaces).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/change-list", ctrl.GetOracleChanges).Methods("GET")

	router.HandleFunc("/hosts/technologies/oracle/exadata", ctrl.SearchOracleExadata).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/total-memory-size", ctrl.GetTotalOracleExadataMemorySizeStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/total-cpu", ctrl.GetTotalOracleExadataCPUStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/average-storage-usage", ctrl.GetAverageOracleExadataStorageUsageStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/storage-error-count-status", ctrl.GetOracleExadataStorageErrorCountStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/patch-status", ctrl.GetOracleExadataPatchStatusStats).Methods("GET")

	// ORACLE CONTRACTS
	router.HandleFunc("/contracts/oracle/database", ctrl.AddOracleDatabaseContract).Methods("POST")
	router.HandleFunc("/contracts/oracle/database", ctrl.UpdateOracleDatabaseContract).Methods("PUT")
	router.HandleFunc("/contracts/oracle/database", ctrl.GetOracleDatabaseContracts).Methods("GET")
	router.HandleFunc("/contracts/oracle/database/{id}", ctrl.DeleteOracleDatabaseContract).Methods("DELETE")

	router.HandleFunc("/contracts/oracle/database/{id}/hosts", ctrl.AddHostToOracleDatabaseContract).Methods("POST")
	router.HandleFunc("/contracts/oracle/database/{id}/hosts/{hostname}", ctrl.DeleteHostFromOracleDatabaseContract).Methods("DELETE")

	// SQL SERVER CONTRACTS
	router.HandleFunc("/contracts/microsoft/database", ctrl.AddSqlServerDatabaseContract).Methods("POST")
	router.HandleFunc("/contracts/microsoft/database", ctrl.UpdateSqlServerDatabaseContract).Methods("PUT")
	router.HandleFunc("/contracts/microsoft/database", ctrl.GetSqlServerDatabaseContracts).Methods("GET")
	router.HandleFunc("/contracts/microsoft/database/{id}", ctrl.DeleteSqlServerDatabaseContract).Methods("DELETE")

	// MYSQL
	router.HandleFunc("/hosts/technologies/mysql/databases", ctrl.SearchMySQLInstances).Methods("GET")
	router.HandleFunc("/hosts/{hostname}/technologies/mysql/databases/{dbname}/ignored/{ignored}", ctrl.UpdateMySqlLicenseIgnoredField).Methods("PUT")

	// MYSQL CONTRACTS
	router.HandleFunc("/contracts/mysql/database", ctrl.AddMySQLContract).Methods("POST")
	router.HandleFunc("/contracts/mysql/database/{id}", ctrl.UpdateMySQLContract).Methods("PUT")
	router.HandleFunc("/contracts/mysql/database", ctrl.GetMySQLContracts).Methods("GET")
	router.HandleFunc("/contracts/mysql/database/{id}", ctrl.DeleteMySQLContract).Methods("DELETE")

	// SQL SERVER
	router.HandleFunc("/hosts/technologies/microsoft/databases", ctrl.SearchSqlServerInstances).Methods("GET")
	router.HandleFunc("/hosts/{hostname}/technologies/microsoft/databases/{dbname}/ignored/{ignored}", ctrl.UpdateSqlServerLicenseIgnoredField).Methods("PUT")

	// POSTGRESQL
	router.HandleFunc("/hosts/technologies/postgresql/databases", ctrl.SearchPostgreSqlInstances).Methods("GET")

	// ALERTS
	router.HandleFunc("/alerts", ctrl.SearchAlerts).Methods("GET")
	router.HandleFunc("/alerts/ack", ctrl.AckAlerts).Methods("POST")

	router.HandleFunc("/database/connection/status", ctrl.GetDatabaseConnectionStatus).Methods("GET")

	ctrl.setupSettingsRoutes(router.PathPrefix("/settings").Subrouter())
	ctrl.setupFrontendAPIRoutes(router.PathPrefix("/frontend").Subrouter())
}

func (ctrl *APIController) setupSettingsRoutes(router *mux.Router) {
	router.HandleFunc("/default-database-tag-choices", ctrl.GetDefaultDatabaseTags).Methods("GET")
	router.HandleFunc("/features", ctrl.GetErcoleFeatures).Methods("GET")
	router.HandleFunc("/technologies", ctrl.GetTechnologyList).Methods("GET")
	router.HandleFunc("/oracle/database/license-types", ctrl.GetOracleDatabaseLicenseTypes).Methods("GET")
	router.HandleFunc("/oracle/database/license-types/{id}", ctrl.DeleteOracleDatabaseLicenseType).Methods("DELETE")
	router.HandleFunc("/oracle/database/license-types", ctrl.AddOracleDatabaseLicenseType).Methods("POST")
	router.HandleFunc("/oracle/database/license-types/{id}", ctrl.UpdateOracleDatabaseLicenseType).Methods("PUT")
	router.HandleFunc("/microsoft/database/license-types", ctrl.GetSqlServerDatabaseLicenseTypes).Methods("GET")
	router.HandleFunc("/mysql/database/license-types", ctrl.GetMySqlLicenseTypes).Methods("GET")
}

func (ctrl *APIController) setupFrontendAPIRoutes(router *mux.Router) {
	router.HandleFunc("/dashboard", ctrl.GetInfoForFrontendDashboard).Methods("GET")
}
