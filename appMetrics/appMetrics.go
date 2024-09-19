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

package appmetrics

import (
	"azure_app_exporter/logging"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	TokenSeconds = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "azure_api_token_update_duration_seconds",
		Help: "How many seconds it takes to update the Azure API token.",
	})
	TokenFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "azure_api_token_update_failures",
		Help: "How many times updating the Azure API token has failed.",
	})
	ApplicationsSeconds = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "azure_applications_update_duration_seconds",
		Help: "How many seconds it takes to update the in-memory cache of Azure applications.",
	})
	ApplicationsFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "azure_applications_update_failures",
		Help: "How many times updating the cached Azure applications has failed.",
	})

	ApplicationPasswordSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "azure_application_password_remaining_seconds",
		Help: "Seconds remaining until the password credential expires.",
	}, []string{"id", "app_id", "app_display_name", "password_key_id", "password_display_name", "password_end_date_time"})
)

func init() {
	if err := prometheus.Register(TokenSeconds); err != nil {
		logging.Fatal(err)
	}
	if err := prometheus.Register(TokenFailures); err != nil {
		logging.Fatal(err)
	}
	if err := prometheus.Register(ApplicationsSeconds); err != nil {
		logging.Fatal(err)
	}
	if err := prometheus.Register(ApplicationsFailures); err != nil {
		logging.Fatal(err)
	}
	if err := prometheus.Register(ApplicationPasswordSeconds); err != nil {
		logging.Fatal(err)
	}
}
