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
	"crypto/tls"
	"fmt"
	"reflect"
)

type (
	CipherSuite     uint16
	ProtocolVersion int
)

var cipherSuiteValue = map[string]CipherSuite{
	"TLS13_AES_256_GCM_SHA384":                      CipherSuite(tls.TLS_AES_256_GCM_SHA384),
	"TLS13_AES_128_GCM_SHA256":                      CipherSuite(tls.TLS_AES_128_GCM_SHA256),
	"TLS13_CHACHA20_POLY1305_SHA256":                CipherSuite(tls.TLS_CHACHA20_POLY1305_SHA256),
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384),
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256),
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256),
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384),
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256),
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   CipherSuite(tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256),
}

var cipherSuiteName = map[CipherSuite]string{
	CipherSuite(tls.TLS_AES_256_GCM_SHA384):                        "TLS13_AES_256_GCM_SHA384",
	CipherSuite(tls.TLS_AES_128_GCM_SHA256):                        "TLS13_AES_128_GCM_SHA256",
	CipherSuite(tls.TLS_CHACHA20_POLY1305_SHA256):                  "TLS13_CHACHA20_POLY1305_SHA256",
	CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384):       "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256):       "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	CipherSuite(tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256): "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
	CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384):         "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	CipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256):         "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	CipherSuite(tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256):   "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
}

var protocolVersionValue = map[string]ProtocolVersion{
	"TLS13": ProtocolVersion(tls.VersionTLS13),
	"TLS12": ProtocolVersion(tls.VersionTLS12),
}

var protocolVersionName = map[ProtocolVersion]string{
	ProtocolVersion(tls.VersionTLS13): "TLS13",
	ProtocolVersion(tls.VersionTLS12): "TLS12",
}

func (c CipherSuite) String() string {
	return cipherSuiteName[c]
}

func (p ProtocolVersion) String() string {
	return protocolVersionName[p]
}

func (c CipherSuite) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (p ProtocolVersion) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (c *CipherSuite) UnmarshalText(bytes []byte) error {
	name := string(bytes)

	if cipherSuite, ok := cipherSuiteValue[name]; ok {
		*c = cipherSuite
		return nil
	}

	return fmt.Errorf("invalid cipher suite %s, expected one of %v", name, reflect.ValueOf(cipherSuiteValue).MapKeys())
}

func (p *ProtocolVersion) UnmarshalText(bytes []byte) error {
	name := string(bytes)

	if protocolVersion, ok := protocolVersionValue[name]; ok {
		*p = protocolVersion
		return nil
	}

	return fmt.Errorf("invalid protocol version %s, expected one of %v", name, reflect.ValueOf(protocolVersionValue).MapKeys())
}
