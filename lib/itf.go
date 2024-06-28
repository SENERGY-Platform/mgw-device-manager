package lib

import (
	"context"
	srv_info_lib "github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib"
	"github.com/SENERGY-Platform/mgw-device-manager/lib/model"
)

type Api interface {
	baseApi
	clientApi
	userApi
}

type ClientApi interface {
	baseApi
	clientApi
}

type UserApi interface {
	baseApi
	userApi
}

type baseApi interface {
	GetDevice(ctx context.Context, id string) (model.Device, error)
	GetDevices(ctx context.Context, filter model.DevicesFilter) (map[string]model.Device, error)
	DeleteDevice(ctx context.Context, id string) error
	srv_info_lib.Api
}

type clientApi interface {
	AddDevice(ctx context.Context, deviceBase model.DeviceBase) error
	UpdateDevice(ctx context.Context, deviceBase model.DeviceBase) error
}

type userApi interface {
	UpdateDeviceUserData(ctx context.Context, id string, userDataBase model.DeviceUserDataBase) error
}
