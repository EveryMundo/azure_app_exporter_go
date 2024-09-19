/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package applications

import (
	"azure_app_exporter/logging"
	"context"
	"fmt"
	"time"

	appmetrics "azure_app_exporter/appMetrics"
	datatypes "azure_app_exporter/azure/applications/dataTypes"
	globalstate "azure_app_exporter/globalState"
)

// https://learn.microsoft.com/en-us/graph/query-parameters
// https://learn.microsoft.com/en-us/graph/api/application-list?view=graph-rest-1.0
func AzureApplicationsUpdater() {
	// This func is spawned in a thread simultaneously with another thread
	// responsible for updating the api token, so we should wait for it to finish
	for globalstate.AzureApiToken.Value == "" {
		logging.Warn("azure api token not yet acquired, sleeping 5 seconds")
		time.Sleep(5 * time.Second)
	}

	httpClient := globalstate.HttpClient.Clone()

	getApplications := func(url string) (datatypes.AzureApplications, error) {
		logging.Debugf("calling with bearer token: %s", url)

		globalstate.AzureApiToken.RwLock.RLock()
		defer globalstate.AzureApiToken.RwLock.RUnlock()

		var response datatypes.AzureApplications
		err := httpClient.
			BaseURL(url).
			Bearer(globalstate.AzureApiToken.Value).
			ToJSON(&response).
			Fetch(context.Background())

		return response, err
	}

	inner := func() error {
		response, err := getApplications(
			fmt.Sprintf(
				"%s?$top=%d&$select=id,appId,displayName,createdDateTime,passwordCredentials",
				globalstate.Settings.Applications.Url,
				globalstate.Settings.Applications.ResultsPerPage,
			),
		)
		if err != nil {
			return err
		}

		for response.NextLink != nil {
			nextResponse, err := getApplications(*response.NextLink)
			if err != nil {
				return err
			}

			response.NextLink = nextResponse.NextLink
			response.Value = append(response.Value, nextResponse.Value...)
		}

		globalstate.Applications.RwLock.Lock()
		defer globalstate.Applications.RwLock.Unlock()

		for k := range globalstate.Applications.Value {
			delete(globalstate.Applications.Value, k)
		}

		for _, application := range response.Value {
			globalstate.Applications.Value[application.Id] = application
		}

		logging.Debugf("cached %d applications", len(globalstate.Applications.Value))

		return nil
	}

	for {
		start := time.Now()

		if err := inner(); err == nil {
			elapsed := time.Since(start)
			appmetrics.ApplicationsSeconds.Observe(elapsed.Seconds())
			logging.Infof("updated azure applications in %s, next update after %s", elapsed, globalstate.Settings.Applications.CacheRefreshInterval)
		} else {
			logging.Errorf("failed updating azure applications -> %s, new attempt after %s", err, globalstate.Settings.Applications.CacheRefreshInterval)
			appmetrics.ApplicationsFailures.Inc()
		}

		time.Sleep(globalstate.Settings.Applications.CacheRefreshInterval.Duration)
	}
}
