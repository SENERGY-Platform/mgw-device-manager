package handler

import (
	"context"
	"database/sql/driver"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"time"
)

type DevicesHandler interface {
	Put(ctx context.Context, deviceData lib_model.DeviceData) error
	Get(ctx context.Context, id string) (lib_model.Device, error)
	GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error)
	SetUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error
	SetStates(ctx context.Context, ref string, state lib_model.DeviceState) error
	Delete(ctx context.Context, id string) error
}

type DevicesStorageHandler interface {
	BeginTransaction(ctx context.Context) (driver.Tx, error)
	Create(ctx context.Context, tx driver.Tx, device lib_model.DeviceBase) error
	Read(ctx context.Context, id string) (lib_model.Device, error)
	ReadAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error)
	Update(ctx context.Context, tx driver.Tx, deviceBase lib_model.DeviceBase) error
	UpdateUserData(ctx context.Context, tx driver.Tx, id string, userData lib_model.DeviceUserData) error
	UpdateStates(ctx context.Context, tx driver.Tx, ref string, state lib_model.DeviceState, timestamp time.Time) error
	Delete(ctx context.Context, tx driver.Tx, id string) error
}

type MqttClient interface {
	Subscribe(topic string, qos byte, messageHandler func(m Message)) error
	Publish(topic string, qos byte, retained bool, payload any) error
}

type Message interface {
	Topic() string
	Payload() []byte
}

type MessageRelayHandler interface {
	Put(m Message) error
}

type MessageHandler func(m Message) error
