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
	"net/http"

	datatypes "azure_app_exporter/azure/applications/dataTypes"
	fromswaggerui "azure_app_exporter/fromSwaggerUi"
	globalstate "azure_app_exporter/globalState"

	"github.com/labstack/echo/v4"
)

// @summary Show all Azure applications cached in the exporter (truncated in Swagger UI to 50 entries)
// @description Show all Azure applications cached in the exporter (truncated in Swagger UI to 50 entries)
// @description
// @description Call this endpoint outside Swagger UI to see full response
// @tags applications
// @produce json
// @success 200 {object} map[string]datatypes.AzureApplication
// @router /api/apps [get]
func AllApplications(c echo.Context) error {
	globalstate.Applications.RwLock.RLock()
	defer globalstate.Applications.RwLock.RUnlock()

	if _, fromUi := c.Request().Header[fromswaggerui.HeaderName]; fromUi {
		i := 0
		truncated := make(map[string]datatypes.AzureApplication, 50)

		for id, application := range globalstate.Applications.Value {
			if i >= 50 {
				break
			}
			truncated[id] = application
			i++
		}

		return c.JSON(http.StatusOK, truncated)
	}

	return c.JSON(http.StatusOK, globalstate.Applications.Value)
}

// @summary Show Azure application by ID
// @description Show Azure application by ID
// @tags applications
// @param id path string true "ID of Azure application to lookup"
// @produce json
// @success 200 {object} datatypes.AzureApplication
// @router /api/apps/{id} [get]
func ApplicationById(c echo.Context) error {
	globalstate.Applications.RwLock.RLock()
	defer globalstate.Applications.RwLock.RUnlock()

	if application, ok := globalstate.Applications.Value[c.Param("id")]; ok {
		return c.JSON(http.StatusOK, application)
	}

	return c.NoContent(http.StatusNotFound)
}
