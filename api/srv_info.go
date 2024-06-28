package api

import (
	"context"
	srv_info_lib "github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib"
)

func (a *Api) GetSrvInfo(_ context.Context) srv_info_lib.SrvInfo {
	return a.srvInfoHdl.GetInfo()
}
