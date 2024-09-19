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
	appmetrics "azure_app_exporter/appMetrics"
	datatypes "azure_app_exporter/azure/applications/dataTypes"
	globalstate "azure_app_exporter/globalState"
)

func derefOrDefault(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func derefOrDefaultUtcTime(u *datatypes.UtcTime) string {
	if u != nil {
		return u.String()
	}

	return ""
}

func UpdateApplicationsMetrics() {
	globalstate.Applications.RwLock.RLock()
	defer globalstate.Applications.RwLock.RUnlock()

	for id, application := range globalstate.Applications.Value {
		for _, password := range application.PasswordCredentials {
			appmetrics.ApplicationPasswordSeconds.WithLabelValues(
				id,
				application.AppId,
				derefOrDefault(application.DisplayName),
				password.KeyId,
				derefOrDefault(password.DisplayName),
				derefOrDefaultUtcTime(password.EndDateTime),
			).
				Set(password.RemainingSeconds())
		}
	}
}
