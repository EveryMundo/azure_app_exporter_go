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

package logging

import (
	"os"

	"github.com/labstack/gommon/log"
)

func init() {
	// Format specifiers https://echo.labstack.com/docs/customization#log-header
	log.SetHeader("${time_rfc3339} ${level} ${short_file}:${line}")
	log.SetLevel(log.DEBUG)
	log.SetOutput(os.Stderr)
}

var (
	Debug  = log.Debug
	Debugf = log.Debugf
	Info   = log.Info
	Infof  = log.Infof
	Warn   = log.Warn
	Warnf  = log.Warnf
	Error  = log.Error
	Errorf = log.Errorf
	Fatal  = log.Fatal
	Fatalf = log.Fatalf
)
