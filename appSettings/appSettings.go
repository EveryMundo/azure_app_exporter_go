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

package appsettings

import (
	"azure_app_exporter/logging"
	"crypto/tls"
	"os"
	"sort"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Settings struct {
	Credentials  Credentials  `toml:"credentials"  json:"credentials"  validate:"required" extensions:"x-order=1"`
	Metrics      Metrics      `toml:"metrics"      json:"metrics"                          extensions:"x-order=2"`
	Applications Applications `toml:"applications" json:"applications"                     extensions:"x-order=3"`
	Web          Web          `toml:"web"          json:"web"                              extensions:"x-order=4"`
	OpenApi      OpenApi      `toml:"openapi"      json:"openapi"                          extensions:"x-order=5"`
	Tls          Tls          `toml:"tls"          json:"tls"                              extensions:"x-order=6"`
	Debug        Debug        `toml:"debug"        json:"debug"                            extensions:"x-order=7"`
}

type Credentials struct {
	TenantId     string       `toml:"tenant_id"     json:"tenant_id"     extensions:"x-order=1"`
	ClientId     string       `toml:"client_id"     json:"client_id"     extensions:"x-order=2"`
	ClientSecret ClientSecret `toml:"client_secret" json:"client_secret" extensions:"x-order=3"`
}

type Metrics struct {
	PruneInterval               *Duration `toml:"prune_interval"                 json:"prune_interval"                 swaggertype:"string" example:"30m" extensions:"x-order=1,x-nullable"`
	ExpandUnsupportedUrlMetrics bool      `toml:"expand_unsupported_url_metrics" json:"expand_unsupported_url_metrics"                                    extensions:"x-order=2"`
}

type Applications struct {
	Enabled              bool     `toml:"enabled"                json:"enabled"                extensions:"x-order=1"`
	CacheRefreshInterval Duration `toml:"cache_refresh_interval" json:"cache_refresh_interval" extensions:"x-order=2" swaggertype:"string" example:"15m"`
	Url                  string   `toml:"url"                    json:"url"                    extensions:"x-order=2"`
	ResultsPerPage       uint16   `toml:"results_per_page"       json:"results_per_page"       extensions:"x-order=3"                                    minimum:"1" maximum:"999"`
}

type Web struct {
	ListenAddress string  `toml:"listen_address" json:"listen_address" extensions:"x-order=1"`
	CertFile      *string `toml:"cert_file"      json:"cert_file"      extensions:"x-order=2,x-nullable"`
	KeyFile       *string `toml:"key_file"       json:"key_file"       extensions:"x-order=3,x-nullable"`
}

type OpenApi struct {
	Enabled      bool   `toml:"enabled"        json:"enabled"        extensions:"x-order=1"`
	DocsUrl      string `toml:"docs_url"       json:"docs_url"       extensions:"x-order=2"`
	SwaggerUiUrl string `toml:"swagger_ui_url" json:"swagger_ui_url" extensions:"x-order=3"`
}

type Tls struct {
	CipherSuites     []CipherSuite     `toml:"cipher_suites"     json:"cipher_suites"     swaggertype:"array,string" example:"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256" extensions:"x-order=1"`
	ProtocolVersions []ProtocolVersion `toml:"protocol_versions" json:"protocol_versions" swaggertype:"array,string" example:"TLS13"                                       extensions:"x-order=2"`
	// key_exchange_groups is not supported in Go crypto/tls
}

func (t Tls) ToCipherSuites() []uint16 {
	cipherSuites := make([]uint16, 0, len(t.CipherSuites))

	for _, cipherSuite := range t.CipherSuites {
		cipherSuites = append(cipherSuites, uint16(cipherSuite))
	}

	return cipherSuites
}

type Debug struct {
	NoVerifyTls bool `toml:"no_verify_tls" json:"no_verify_tls"`
}

func Parse() Settings {
	settingsEnvVar := "AZURE_APP_EXPORTER_SETTINGS_PATH"
	settingsPath := "/etc/azure_app_exporter/settings.toml"

	if path, ok := os.LookupEnv(settingsEnvVar); ok {
		settingsPath = path
	} else {
		logging.Warnf("no %s env var set, defaulting to %s", settingsEnvVar, settingsPath)
	}

	settingsContents, err := os.ReadFile(settingsPath)
	if err != nil {
		logging.Fatalf("failed reading %s -> %s", settingsPath, err)
	}

	settings := Settings{
		Applications: Applications{
			Enabled:              true,
			CacheRefreshInterval: Duration{15 * time.Minute},
			Url:                  "https://graph.microsoft.com/v1.0/applications",
			ResultsPerPage:       999,
		},
		Web: Web{
			ListenAddress: "0.0.0.0:9081",
		},
		OpenApi: OpenApi{
			Enabled:      true,
			DocsUrl:      "/openapi.json",
			SwaggerUiUrl: "/swagger",
		},
		Tls: Tls{
			CipherSuites: []CipherSuite{
				CipherSuite(tls.TLS_AES_256_GCM_SHA384),
				CipherSuite(tls.TLS_AES_128_GCM_SHA256),
				CipherSuite(tls.TLS_CHACHA20_POLY1305_SHA256),
				CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384),
				CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256),
				CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256),
				CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384),
				CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256),
				CipherSuite(tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256),
			},
			ProtocolVersions: []ProtocolVersion{
				ProtocolVersion(tls.VersionTLS13),
				ProtocolVersion(tls.VersionTLS12),
			},
		},
	}

	if err := toml.Unmarshal(settingsContents, &settings); err != nil {
		logging.Fatalf("failed parsing %s -> %s", settingsPath, err)
	}

	sort.Slice(settings.Tls.ProtocolVersions, func(i, j int) bool {
		return settings.Tls.ProtocolVersions[i] < settings.Tls.ProtocolVersions[j]
	})

	validate(settings)

	return settings
}

func validate(s Settings) {
	if s.Applications.ResultsPerPage < 1 || s.Applications.ResultsPerPage > 999 {
		logging.Fatalf("settings value %d not in range 1..=999", s.Applications.ResultsPerPage)
	}

	if len(s.Tls.ProtocolVersions) < 1 {
		logging.Fatal("tls protocol versions cannot be empty")
	}

	checkUrl := func(url string) {
		if url == "" || url == "/" {
			logging.Fatalf("url %s cannot be empty or \"/\"", url)
		}
	}

	checkUrl(s.OpenApi.DocsUrl)
	checkUrl(s.OpenApi.SwaggerUiUrl)

	verifyCredentialPresent := func(credential string) {
		if credential == "" || credential == "..." {
			logging.Fatalf("empty credential found in settings.toml: %s", credential)
		}
	}

	verifyCredentialPresent(s.Credentials.TenantId)
	verifyCredentialPresent(s.Credentials.ClientId)
	verifyCredentialPresent(string(s.Credentials.ClientSecret))
}
