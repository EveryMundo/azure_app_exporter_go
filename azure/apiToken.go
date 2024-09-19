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

package azure

import (
	"azure_app_exporter/logging"
	"context"
	"fmt"
	"net/url"
	"time"

	appmetrics "azure_app_exporter/appMetrics"
	globalstate "azure_app_exporter/globalState"
)

type authToken struct {
	ExpiresIn   uint64 `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

// https://learn.microsoft.com/en-us/graph/auth-v2-service#4-request-an-access-token
func AzureApiTokenUpdater() {
	httpClient := globalstate.HttpClient.Clone()

	inner := func() (time.Duration, error) {
		requestUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", globalstate.Settings.Credentials.TenantId)
		logging.Debugf("calling with client id and secret: %s", requestUrl)

		var response authToken
		if err := httpClient.
			BaseURL(requestUrl).
			Post().
			BodyForm(url.Values{
				"grant_type":    {"client_credentials"},
				"scope":         {"https://graph.microsoft.com/.default"},
				"client_id":     {globalstate.Settings.Credentials.ClientId},
				"client_secret": {string(globalstate.Settings.Credentials.ClientSecret)},
			}).
			ToJSON(&response).
			Fetch(context.Background()); err != nil {
			return 0, err
		}

		globalstate.AzureApiToken.RwLock.Lock()
		defer globalstate.AzureApiToken.RwLock.Unlock()

		globalstate.AzureApiToken.Value = response.AccessToken

		return time.Duration(response.ExpiresIn) * time.Second, nil
	}

	for {
		start := time.Now()

		sleepDuration := 30 * time.Second

		if duration, err := inner(); err == nil {
			elapsed := time.Since(start)
			sleepDuration = time.Duration(duration.Seconds()*0.9) * time.Second // Sleep for 90% of the token's validity duration
			logging.Infof("updated azure api token in %s, next update after %s", elapsed, sleepDuration)
			appmetrics.TokenSeconds.Observe(elapsed.Seconds())
		} else {
			logging.Errorf("failed updating api token -> %s, new attempt after %s", err, sleepDuration)
			appmetrics.TokenFailures.Inc()
		}

		time.Sleep(sleepDuration)
	}
}
