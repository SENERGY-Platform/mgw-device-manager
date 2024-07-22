package message_hdl

import (
	"context"
	"encoding/json"
	"errors"
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"reflect"
	"testing"
)

func TestHandler_HandleMessage(t *testing.T) {
	util.InitLogger(sb_util.LoggerConfig{Terminal: true, Level: 4})
	t.Run("set device", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{
			Devices: make(map[string]lib_model.DeviceDataBase),
		}
		h := Handler{devicesHdl: mockDHdl}
		a := lib_model.DeviceDataBase{
			ID:    "123",
			Ref:   "test",
			Name:  "test",
			State: lib_model.Online,
			Type:  "test2",
			Attributes: []lib_model.DeviceAttribute{
				{
					Key:   "a",
					Value: "b",
				},
			},
		}
		p, err := json.Marshal(lib_model.DeviceMessage{
			Method:   lib_model.Set,
			DeviceID: "123",
			Data: &lib_model.DeviceMessageData{
				Name:  "test",
				State: lib_model.Online,
				Type:  "test2",
				Attributes: []lib_model.DeviceAttribute{
					{
						Key:   "a",
						Value: "b",
					},
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
		b, ok := mockDHdl.Devices["123"]
		if !ok {
			t.Error("not in map")
		}
		if !reflect.DeepEqual(a, b) {
			t.Error("got", b, "expected", a)
		}
		if mockDHdl.PutC != 1 {
			t.Error("missing call")
		}
		t.Run("no device data", func(t *testing.T) {
			mockDHdl := &mockDeviceHdl{}
			h := Handler{devicesHdl: mockDHdl}
			p2, err := json.Marshal(lib_model.DeviceMessage{
				Method:   lib_model.Set,
				DeviceID: "123",
			})
			if err != nil {
				t.Fatal(err)
			}
			h.HandleMessage(&mockMessage{
				topic:   "device-manager/device/test",
				payload: p2,
			})
		})
		t.Run("error", func(t *testing.T) {
			mockDHdl := &mockDeviceHdl{PutErr: errors.New("test")}
			h := Handler{devicesHdl: mockDHdl}
			h.HandleMessage(&mockMessage{
				topic:   "device-manager/device/test",
				payload: p,
			})
		})
	})
	t.Run("delete device", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{}
		h := Handler{devicesHdl: mockDHdl}
		p, err := json.Marshal(lib_model.DeviceMessage{
			Method:   lib_model.Delete,
			DeviceID: "123",
		})
		if err != nil {
			t.Fatal(err)
		}
		h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
		if mockDHdl.DeleteC != 1 {
			t.Error("missing call")
		}
		t.Run("error", func(t *testing.T) {
			mockDHdl := &mockDeviceHdl{DeleteErr: errors.New("test")}
			h := Handler{devicesHdl: mockDHdl}
			h.HandleMessage(&mockMessage{
				topic:   "device-manager/device/test",
				payload: p,
			})
		})
	})
	t.Run("set states", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{
			States: make(map[string]lib_model.DeviceState),
		}
		h := Handler{devicesHdl: mockDHdl}
		h.HandleMessage(&mockMessage{
			topic: "device-manager/device/test/lw",
		})
		s, ok := mockDHdl.States["test"]
		if !ok {
			t.Error("not in map")
		}
		if s != lib_model.Offline {
			t.Error("got", s, "expected", lib_model.Offline)
		}
		t.Run("error", func(t *testing.T) {
			mockDHdl := &mockDeviceHdl{SetStatesErr: errors.New("test")}
			h := Handler{devicesHdl: mockDHdl}
			h.HandleMessage(&mockMessage{
				topic: "device-manager/device/test/lw",
			})
		})
	})
	t.Run("unknown method", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{}
		h := Handler{devicesHdl: mockDHdl}
		p, err := json.Marshal(lib_model.DeviceMessage{
			Method: "test",
		})
		if err != nil {
			t.Fatal(err)
		}
		h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
	})
	t.Run("parse topic error", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{}
		h := Handler{devicesHdl: mockDHdl}
		h.HandleMessage(&mockMessage{
			topic: "test",
		})
	})
}

type mockDeviceHdl struct {
	Devices      map[string]lib_model.DeviceDataBase
	States       map[string]lib_model.DeviceState
	PutErr       error
	SetStatesErr error
	DeleteErr    error
	PutC         int
	SetStatesC   int
	DeleteC      int
}

func (m *mockDeviceHdl) Put(ctx context.Context, deviceData lib_model.DeviceDataBase) error {
	m.PutC++
	if m.PutErr != nil {
		return m.PutErr
	}
	m.Devices[deviceData.ID] = deviceData
	return nil
}

func (m *mockDeviceHdl) Get(ctx context.Context, id string) (lib_model.DeviceBase, error) {
	panic("not implemented")
}

func (m *mockDeviceHdl) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.DeviceBase, error) {
	panic("not implemented")
}

func (m *mockDeviceHdl) SetUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	panic("not implemented")
}

func (m *mockDeviceHdl) SetStates(ctx context.Context, ref string, state lib_model.DeviceState) error {
	m.SetStatesC++
	if m.SetStatesErr != nil {
		return m.SetStatesErr
	}
	m.States[ref] = state
	return nil
}

func (m *mockDeviceHdl) Delete(ctx context.Context, id string) error {
	m.DeleteC++
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	return nil
}

type mockMessage struct {
	topic   string
	payload []byte
}

func (m *mockMessage) Topic() string {
	return m.topic
}

func (m *mockMessage) Payload() []byte {
	return m.payload
}
