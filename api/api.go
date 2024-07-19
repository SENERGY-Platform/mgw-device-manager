package api

import (
	srv_info_hdl "github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
)

type Api struct {
	devicesHdl handler.DevicesHandler
	srvInfoHdl srv_info_hdl.SrvInfoHandler
}

func New(devicesHdl handler.DevicesHandler, srvInfoHdl srv_info_hdl.SrvInfoHandler) *Api {
	return &Api{
		devicesHdl: devicesHdl,
		srvInfoHdl: srvInfoHdl,
	}
}
