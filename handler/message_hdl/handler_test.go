package message_hdl

import (
	"context"
	"encoding/json"
	"errors"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"reflect"
	"testing"
)

func TestHandler_HandleMessage(t *testing.T) {
	t.Run("set device", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{
			Devices: make(map[string]lib_model.DeviceData),
		}
		h := Handler{devicesHdl: mockDHdl}
		a := lib_model.DeviceData{
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
		err = h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
		if err != nil {
			t.Error(err)
		}
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
			err = h.HandleMessage(&mockMessage{
				topic:   "device-manager/device/test",
				payload: p2,
			})
			if err == nil {
				t.Error("expected error")
			}
		})
		t.Run("error", func(t *testing.T) {
			mockDHdl := &mockDeviceHdl{PutErr: lib_model.NewInvalidInputError(errors.New("test"))}
			h := Handler{devicesHdl: mockDHdl}
			err = h.HandleMessage(&mockMessage{
				topic:   "device-manager/device/test",
				payload: p,
			})
			if err == nil {
				t.Error("expected error")
			}
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
		_ = h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
		if mockDHdl.DeleteC != 1 {
			t.Error("missing call")
		}
	})
	t.Run("set states", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{
			States: make(map[string]lib_model.DeviceState),
		}
		h := Handler{devicesHdl: mockDHdl}
		_ = h.HandleMessage(&mockMessage{
			topic: "device-manager/device/test/lw",
		})
		s, ok := mockDHdl.States["test"]
		if !ok {
			t.Error("not in map")
		}
		if s != lib_model.Offline {
			t.Error("got", s, "expected", lib_model.Offline)
		}
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
		err = h.HandleMessage(&mockMessage{
			topic:   "device-manager/device/test",
			payload: p,
		})
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("parse topic error", func(t *testing.T) {
		mockDHdl := &mockDeviceHdl{}
		h := Handler{devicesHdl: mockDHdl}
		err := h.HandleMessage(&mockMessage{
			topic: "test",
		})
		if err == nil {
			t.Error("expected error")
		}
	})
}

type mockDeviceHdl struct {
	Devices      map[string]lib_model.DeviceData
	States       map[string]lib_model.DeviceState
	PutErr       error
	SetStatesErr error
	DeleteErr    error
	PutC         int
	SetStatesC   int
	DeleteC      int
}

func (m *mockDeviceHdl) Put(ctx context.Context, deviceData lib_model.DeviceData) error {
	m.PutC++
	if m.PutErr != nil {
		return m.PutErr
	}
	m.Devices[deviceData.ID] = deviceData
	return nil
}

func (m *mockDeviceHdl) Get(ctx context.Context, id string) (lib_model.Device, error) {
	panic("not implemented")
}

func (m *mockDeviceHdl) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
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
