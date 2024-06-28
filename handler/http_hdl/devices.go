package http_hdl

import (
	"github.com/SENERGY-Platform/mgw-device-manager/lib"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const devIdParam = "d"

type devicesQuery struct {
	IDs   string `form:"ids"`
	State string `form:"state"`
	Type  string `form:"type"`
	Ref   string `form:"ref"`
}

func getDevicesH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var query devicesQuery
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		devices, err := a.GetDevices(gc.Request.Context(), lib_model.DevicesFilter{
			IDs:   parseStringSlice(query.IDs, ","),
			State: query.State,
			Type:  query.Type,
			Ref:   query.Ref,
		})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, devices)
	}
}

func getDeviceH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		device, err := a.GetDevice(gc.Request.Context(), gc.Param(devIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, device)
	}
}

func patchUpdateDeviceUserDataH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var userDataBase lib_model.DeviceUserDataBase
		err := gc.ShouldBindJSON(&userDataBase)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err = a.UpdateDeviceUserData(gc.Request.Context(), gc.Param(devIdParam), userDataBase)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func deleteDeviceH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		err := a.DeleteDevice(gc.Request.Context(), gc.Param(devIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
