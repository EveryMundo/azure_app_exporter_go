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

package globalstate

import (
	"crypto/tls"
	"net/http"
	"sync"

	appsettings "azure_app_exporter/appSettings"
	datatypes "azure_app_exporter/azure/applications/dataTypes"

	"github.com/carlmjohnson/requests"
)

var (
	Settings     = appsettings.Parse()
	HttpClient   = requests.Builder{}
	Applications = struct {
		// map of id -> application
		Value  map[string]datatypes.AzureApplication
		RwLock sync.RWMutex
	}{Value: make(map[string]datatypes.AzureApplication)}
	AzureApiToken struct {
		Value  string
		RwLock sync.RWMutex
	}
)

func init() {
	if Settings.Debug.NoVerifyTls {
		HttpClient.Transport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
	}
}
