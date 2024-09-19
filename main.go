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

package main

import (
	"azure_app_exporter/azure"
	"azure_app_exporter/azure/applications"
	"azure_app_exporter/logging"
	"azure_app_exporter/pages"
	"crypto/tls"
	"net/http"
	"os"
	"strings"

	apisettings "azure_app_exporter/appSettings/api"

	_ "azure_app_exporter/docs"
	fromswaggerui "azure_app_exporter/fromSwaggerUi"
	globalstate "azure_app_exporter/globalState"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Azure app exporter
// @version 0.1.0
// TODO choose license
// @description Expose Prometheus metrics for expiring Azure password credentials
func main() {
	if globalstate.Settings.Debug.NoVerifyTls {
		logging.Warn("flag no_verify_tls is enabled, CERTIFICATES ON FOREIGN API REQUESTS WILL NOT BE VALIDATED!")
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(
		middleware.Recover(),
		middleware.LoggerWithConfig(
			middleware.LoggerConfig{
				// This logger does not support ${level}
				// Format specifiers https://pkg.go.dev/github.com/labstack/echo/v4@v4.11.2/middleware#LoggerConfig
				Format: "${time_rfc3339} request{method=${method} uri=${host}${uri} version=${protocol}}: latency=${latency_human} status=${status}\n",
				Output: os.Stderr,
			},
		),
		echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Subsystem: "azure_app_exporter",
			LabelFuncs: map[string]echoprometheus.LabelValueFunc{
				"url": func(c echo.Context, err error) string { // Replace the default "url" label
					if c.Path() != "" { // The current path is a registered endpoint
						return c.Path()
					} else if globalstate.Settings.Metrics.ExpandUnsupportedUrlMetrics {
						return c.Request().URL.String()
					}
					return "unsupported-url"
				},
			},
			HistogramOptsFunc: func(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
				if strings.HasSuffix(opts.Name, "_duration_seconds") {
					opts.Buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0}
				} else if strings.HasSuffix(opts.Name, "_size_bytes") {
					opts.Buckets = []float64{
						512.0,       // 512 B
						1024.0,      // 1,024 B
						2048.0,      // 2.0 KiB
						5120.0,      // 5.0 KiB
						10240.0,     // 10.0 KiB
						25600.0,     // 25.0 KiB
						51200.0,     // 50.0 KiB
						102400.0,    // 100.0 KiB
						256000.0,    // 250.0 KiB
						512000.0,    // 500.0 KiB
						1048576.0,   // 1024.0 KiB
						2097152.0,   // 2.0 MiB
						5242880.0,   // 5.0 MiB
						10485760.0,  // 10.0 MiB
						26214400.0,  // 25.0 MiB
						52428800.0,  // 50.0 MiB
						104857600.0, // 100.0 MiB
					}
				}
				return opts
			},
		}),
		fromswaggerui.SetSwaggerUiHeader,
	)

	if globalstate.Settings.Applications.Enabled {
		go azure.AzureApiTokenUpdater()
		go applications.AzureApplicationsUpdater()
	}

	if globalstate.Settings.OpenApi.Enabled {
		e.GET(globalstate.Settings.OpenApi.SwaggerUiUrl+"/*", echoSwagger.WrapHandler)
	}

	e.GET("/licenses", pages.Licenses)
	e.GET("/metrics", pages.Metrics)
	e.GET("/api/settings", apisettings.ApiSettings)
	e.GET("/api/apps", applications.AllApplications)
	e.GET("/api/apps/:id", applications.ApplicationById)

	logging.Infof("beginning to serve on %s", globalstate.Settings.Web.ListenAddress)
	logging.Infof("metrics endpoint: %s", globalstate.Settings.Web.ListenAddress+"/metrics")
	logging.Infof("swagger endpoint: %s", globalstate.Settings.Web.ListenAddress+globalstate.Settings.OpenApi.SwaggerUiUrl+"/index.html")

	if globalstate.Settings.Web.CertFile != nil && globalstate.Settings.Web.KeyFile != nil {
		server := http.Server{
			Addr:    globalstate.Settings.Web.ListenAddress,
			Handler: e,
			TLSConfig: &tls.Config{
				CipherSuites: globalstate.Settings.Tls.ToCipherSuites(),
				MinVersion:   uint16(globalstate.Settings.Tls.ProtocolVersions[0]),
				MaxVersion:   uint16(globalstate.Settings.Tls.ProtocolVersions[len(globalstate.Settings.Tls.ProtocolVersions)-1]),
			},
		}

		e.Logger.Fatal(server.ListenAndServeTLS(*globalstate.Settings.Web.CertFile, *globalstate.Settings.Web.KeyFile))
	} else {
		logging.Warn("no cert or key file provided in settings.toml, running server in HTTP mode")
		e.Logger.Fatal(e.Start(globalstate.Settings.Web.ListenAddress))
	}
}
