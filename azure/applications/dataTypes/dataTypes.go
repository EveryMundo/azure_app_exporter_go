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

package datatypes

import (
	"math"
	"time"
)

// https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-1.0#properties
type AzureApplications struct {
	NextLink *string            `json:"@odata.nextLink"`
	Value    []AzureApplication `json:"value"`
}

type AzureApplication struct {
	Id                  string               `json:"id"                  validate:"required" extensions:"x-order=1"`
	AppId               string               `json:"appId"               validate:"required" extensions:"x-order=2"`
	DisplayName         *string              `json:"displayName"                             extensions:"x-order=3,x-nullable"`
	PasswordCredentials []PasswordCredential `json:"passwordCredentials" validate:"required" extensions:"x-order=4"`
}

type PasswordCredential struct {
	KeyId       string   `json:"keyId"       validate:"required" extensions:"x-order=1"`
	DisplayName *string  `json:"displayName"                     extensions:"x-order=2,x-nullable"`
	EndDateTime *UtcTime `json:"endDateTime"                     extensions:"x-order=3,x-nullable" swaggertype:"string" format:"date-time"`
}

// Return the remaining seconds until the password credential expires
// If an end time is not set, return positive infinity
func (p PasswordCredential) RemainingSeconds() float64 {
	if p.EndDateTime == nil {
		return math.Inf(1)
	}

	return time.Until(p.EndDateTime.Time).Seconds()
}
