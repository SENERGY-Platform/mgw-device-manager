/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http_hdl

import (
	"github.com/SENERGY-Platform/mgw-device-manager/lib"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/gin-gonic/gin"
	"sort"
)

func SetRoutes(e *gin.Engine, a lib.Api) {
	e.GET(lib_model.DevicesPath, getDevicesH(a))
	e.GET(lib_model.DevicesPath+"/:"+devIdParam, getDeviceH(a))
	e.PATCH(lib_model.DevicesPath+"/:"+devIdParam, patchUpdateDeviceUserDataH(a))
	e.DELETE(lib_model.DevicesPath+"/:"+devIdParam, deleteDeviceH(a))
	e.GET(lib_model.SrvInfoPath, getSrvInfoH(a))
	e.GET("health-check", getServiceHealthH(a))
}

func GetRoutes(e *gin.Engine) [][2]string {
	routes := e.Routes()
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})
	var rInfo [][2]string
	for _, info := range routes {
		rInfo = append(rInfo, [2]string{info.Method, info.Path})
	}
	return rInfo
}

func GetPathFilter() []string {
	return []string{
		"/health-check",
	}
}
