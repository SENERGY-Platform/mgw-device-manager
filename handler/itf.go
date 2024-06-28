package handler

import (
	"context"
	"database/sql/driver"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
)

type DevicesHandler interface {
	Set(ctx context.Context, deviceBase lib_model.DeviceBase) error
	Add(ctx context.Context, deviceBase lib_model.DeviceBase) error
	Get(ctx context.Context, id string) (lib_model.Device, error)
	GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error)
	Update(ctx context.Context, deviceBase lib_model.DeviceBase) error
	UpdateUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error
	Delete(ctx context.Context, id string) error
}

type DevicesStorageHandler interface {
	BeginTransaction(ctx context.Context) (driver.Tx, error)
	Create(ctx context.Context, tx driver.Tx, device lib_model.Device) error
	Read(ctx context.Context, id string) (lib_model.Device, error)
	ReadAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error)
	Update(ctx context.Context, tx driver.Tx, device lib_model.Device) error
	Delete(ctx context.Context, tx driver.Tx, id string) error
}
