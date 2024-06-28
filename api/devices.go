package api

import (
	"context"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
)

func (a *Api) GetDevice(ctx context.Context, id string) (lib_model.Device, error) {
	return a.devicesHdl.Get(ctx, id)
}

func (a *Api) GetDevices(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	return a.devicesHdl.GetAll(ctx, filter)
}

func (a *Api) DeleteDevice(ctx context.Context, id string) error {
	return a.devicesHdl.Delete(ctx, id)
}

func (a *Api) UpdateDeviceUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	return a.devicesHdl.UpdateUserData(ctx, id, userDataBase)
}
