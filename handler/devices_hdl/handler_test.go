package devices_hdl

import (
	"context"
	"database/sql/driver"
	"errors"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"reflect"
	"testing"
)

var id = "1"
var deviceBase = lib_model.DeviceBase{
	ID:  id,
	Ref: "test",
	DeviceData: lib_model.DeviceData{
		Name:  "test",
		State: lib_model.Online,
		Type:  "test",
		Attributes: []lib_model.DeviceAttribute{
			{
				Key:   "a",
				Value: "b",
			},
		},
	},
}

func TestHandler_Set(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		err := h.Set(context.Background(), deviceBase)
		if err != nil {
			t.Error(err)
		}
		device, ok := stgHdl.devices[id]
		if !ok {
			t.Error("not created")
		}
		if !reflect.DeepEqual(deviceBase, device.DeviceBase) {
			t.Error("expected\n", deviceBase, "got\n", device.DeviceBase)
		}
		if device.Created.IsZero() {
			t.Error("created timestamp is zero")
		}
	})
	t.Run("exist", func(t *testing.T) {
		deviceBase2 := deviceBase
		deviceBase2.Name = "test2"
		if err := h.Set(context.Background(), deviceBase2); err != nil {
			t.Error(err)
		}
		device := stgHdl.devices[id]
		if !reflect.DeepEqual(deviceBase2, device.DeviceBase) {
			t.Error("expected\n", deviceBase2, "got\n", device.DeviceBase)
		}
		if device.Updated.IsZero() {
			t.Error("updated timestamp is zero")
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		err := h.Set(context.Background(), lib_model.DeviceBase{})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_Add(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		err := h.Add(context.Background(), deviceBase)
		if err != nil {
			t.Error(err)
		}
		device, ok := stgHdl.devices[id]
		if !ok {
			t.Error("not created")
		}
		if !reflect.DeepEqual(deviceBase, device.DeviceBase) {
			t.Error("expected\n", deviceBase, "got\n", device.DeviceBase)
		}
		if device.Created.IsZero() {
			t.Error("created timestamp is zero")
		}
	})
	t.Run("exist", func(t *testing.T) {
		err := h.Add(context.Background(), deviceBase)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		err := h.Add(context.Background(), lib_model.DeviceBase{})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_Get(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		_, err := h.Get(context.Background(), "test")
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.Device{DeviceBase: deviceBase}
		device, err := h.Get(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(device.DeviceBase, deviceBase) {
			t.Error("expected\n", deviceBase, "got\n", device.DeviceBase)
		}
	})
}

func TestHandler_GetAll(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("no entries", func(t *testing.T) {
		devices, err := h.GetAll(context.Background(), lib_model.DevicesFilter{})
		if err != nil {
			t.Error(err)
		}
		if len(devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
	t.Run("with entries", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.Device{DeviceBase: deviceBase}
		devices, err := h.GetAll(context.Background(), lib_model.DevicesFilter{})
		if err != nil {
			t.Error(err)
		}
		if len(devices) != 1 {
			t.Error("expected 1 entry")
		}
	})
	t.Run("error", func(t *testing.T) {
		stgHdl.getAllErr = errors.New("test error")
		_, err := h.GetAll(context.Background(), lib_model.DevicesFilter{})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_Update(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		if err := h.Update(context.Background(), deviceBase); err == nil {
			t.Error("expected error")
		}
		if len(stgHdl.devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
	stgHdl.devices[id] = lib_model.Device{DeviceBase: deviceBase}
	t.Run("exists", func(t *testing.T) {
		deviceBase2 := deviceBase
		deviceBase2.Name = "test2"
		if err := h.Update(context.Background(), deviceBase2); err != nil {
			t.Error(err)
		}
		device := stgHdl.devices[id]
		if !reflect.DeepEqual(deviceBase2, device.DeviceBase) {
			t.Error("expected\n", deviceBase2, "got\n", device.DeviceBase)
		}
		if device.Updated.IsZero() {
			t.Error("updated timestamp is zero")
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		err := h.Update(context.Background(), lib_model.DeviceBase{})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_UpdateUserData(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	userDataBase := lib_model.DeviceUserDataBase{
		Name: "test",
		Attributes: []lib_model.DeviceAttribute{
			{
				Key:   "a",
				Value: "b",
			},
		},
	}
	t.Run("does not exist", func(t *testing.T) {
		if err := h.UpdateUserData(context.Background(), id, userDataBase); err == nil {
			t.Error("expected error")
		}
		if len(stgHdl.devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.Device{DeviceBase: deviceBase}
		if err := h.UpdateUserData(context.Background(), id, userDataBase); err != nil {
			t.Error(err)
		}
		device := stgHdl.devices[id]
		if !reflect.DeepEqual(userDataBase, device.UserData.DeviceUserDataBase) {
			t.Error("expected\n", userDataBase, "got\n", device.UserData)
		}
		if device.UserData.Updated.IsZero() {
			t.Error("updated timestamp is zero")
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		err := h.UpdateUserData(context.Background(), id, lib_model.DeviceUserDataBase{Attributes: []lib_model.DeviceAttribute{{Value: "test"}}})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_Delete(t *testing.T) {
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.Device)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		if err := h.Delete(context.Background(), id); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.Device{DeviceBase: deviceBase}
		if err := h.Delete(context.Background(), id); err != nil {
			t.Error(err)
		}
		if len(stgHdl.devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
}

func Test_validateDeviceBase(t *testing.T) {
	dBase := lib_model.DeviceBase{
		ID:  "test",
		Ref: "test",
		DeviceData: lib_model.DeviceData{
			State: lib_model.Online,
			Type:  "test",
			Attributes: []lib_model.DeviceAttribute{
				{
					Key:   "test",
					Value: "test",
				},
			},
		},
	}
	t.Run("valid", func(t *testing.T) {
		if err := validateDeviceBase(dBase); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid ID", func(t *testing.T) {
		idb := dBase
		idb.ID = ""
		if err := validateDeviceBase(idb); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid type", func(t *testing.T) {
		idb := dBase
		idb.Type = ""
		if err := validateDeviceBase(idb); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid state", func(t *testing.T) {
		idb := dBase
		idb.State = "test"
		if err := validateDeviceBase(idb); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid ref", func(t *testing.T) {
		idb := dBase
		idb.Ref = ""
		if err := validateDeviceBase(idb); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid attribute", func(t *testing.T) {
		idb := dBase
		idb.Attributes = append(idb.Attributes, lib_model.DeviceAttribute{})
		if err := validateDeviceBase(idb); err == nil {
			t.Error("expected error")
		}
	})
}

func Test_validateDeviceAttribute(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		if err := validateAttributes([]lib_model.DeviceAttribute{}); err != nil {
			t.Error(err)
		}
	})
	t.Run("valid", func(t *testing.T) {
		if err := validateAttributes([]lib_model.DeviceAttribute{{
			Key:   "a",
			Value: "b",
		}}); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if err := validateAttributes([]lib_model.DeviceAttribute{{Value: "b"}}); err == nil {
			t.Error("expected error")
		}
	})
}

func Test_isValidDeviceState(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run(lib_model.Online, func(t *testing.T) {
			if !isValidDeviceState(lib_model.Online) {
				t.Error("expected true")
			}
		})
		t.Run(lib_model.Offline, func(t *testing.T) {
			if !isValidDeviceState(lib_model.Offline) {
				t.Error("expected true")
			}
		})
		t.Run("empty", func(t *testing.T) {
			if !isValidDeviceState("") {
				t.Error("expected true")
			}
		})
	})
	t.Run("invalid", func(t *testing.T) {
		if isValidDeviceState("test") {
			t.Error("expected false")
		}
	})
}

type stgHdlMock struct {
	devices   map[string]lib_model.Device
	getAllErr error
}

func (m *stgHdlMock) BeginTransaction(_ context.Context) (driver.Tx, error) {
	panic("not implemented")
}

func (m *stgHdlMock) Create(_ context.Context, tx driver.Tx, device lib_model.Device) error {
	if tx != nil {
		panic("not implemented")
	}
	if _, ok := m.devices[device.ID]; ok {
		return errors.New("duplicate device")
	}
	m.devices[device.ID] = device
	return nil
}

func (m *stgHdlMock) Read(_ context.Context, id string) (lib_model.Device, error) {
	device, ok := m.devices[id]
	if !ok {
		return device, lib_model.NewNotFoundError(errors.New("not found"))
	}
	return device, nil
}

func (m *stgHdlMock) ReadAll(_ context.Context, _ lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.devices, nil
}

func (m *stgHdlMock) Update(_ context.Context, tx driver.Tx, device lib_model.Device) error {
	if tx != nil {
		panic("not implemented")
	}
	_, ok := m.devices[device.ID]
	if !ok {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	m.devices[device.ID] = device
	return nil
}

func (m *stgHdlMock) Delete(_ context.Context, tx driver.Tx, id string) error {
	if tx != nil {
		panic("not implemented")
	}
	_, ok := m.devices[id]
	if !ok {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	delete(m.devices, id)
	return nil
}
