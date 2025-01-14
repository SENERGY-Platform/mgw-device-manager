package devices_hdl

import (
	"context"
	"database/sql/driver"
	"errors"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"reflect"
	"testing"
)

var id = "1"
var state = lib_model.Online
var deviceData = lib_model.DeviceData{
	DeviceDataBase: lib_model.DeviceDataBase{
		ID:   id,
		Ref:  "test",
		Name: "test",
		Type: "test",
		Attributes: []lib_model.DeviceAttribute{
			{
				Key:   "a",
				Value: "b",
			},
		},
	},
}

func TestHandler_Put(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.DeviceBase)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		err := h.Put(context.Background(), deviceData.DeviceDataBase, state)
		if err != nil {
			t.Error(err)
		}
		device, ok := stgHdl.devices[id]
		if !ok {
			t.Error("not created")
		}
		if !reflect.DeepEqual(deviceData.DeviceDataBase, device.DeviceDataBase) {
			t.Error("expected\n", deviceData.DeviceDataBase, "got\n", device.DeviceDataBase)
		}
		if device.Created.IsZero() {
			t.Error("created timestamp is zero")
		}
		sItem, ok := h.states[id]
		if !ok {
			t.Error("not in map")
		}
		if sItem.ref != deviceData.Ref {
			t.Error("expected\n", deviceData.Ref, "got\n", sItem.ref)
		}
		if sItem.value != state {
			t.Error("expected\n", state, "got\n", sItem.value)
		}
	})
	t.Run("exist", func(t *testing.T) {
		deviceData2 := deviceData
		deviceData2.Name = "test2"
		if err := h.Put(context.Background(), deviceData2.DeviceDataBase, lib_model.Offline); err != nil {
			t.Error(err)
		}
		device := stgHdl.devices[id]
		if !reflect.DeepEqual(deviceData2.DeviceDataBase, device.DeviceDataBase) {
			t.Error("expected\n", deviceData2.DeviceDataBase, "got\n", device.DeviceDataBase)
		}
		if device.Updated.IsZero() {
			t.Error("updated timestamp is zero")
		}
		sItem, ok := h.states[id]
		if !ok {
			t.Error("not in map")
		}
		if sItem.ref != deviceData.Ref {
			t.Error("expected\n", deviceData.Ref, "got\n", sItem.ref)
		}
		if sItem.value != lib_model.Offline {
			t.Error("expected\n", lib_model.Offline, "got\n", sItem.value)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		t.Run("device data", func(t *testing.T) {
			err := h.Put(context.Background(), lib_model.DeviceDataBase{}, "")
			if err == nil {
				t.Error("expected error")
			}
		})
		t.Run("state", func(t *testing.T) {
			err := h.Put(context.Background(), deviceData.DeviceDataBase, "test")
			if err == nil {
				t.Error("expected error")
			}
		})
	})
}

func TestHandler_Get(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.DeviceBase)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		_, err := h.Get(context.Background(), "test")
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.DeviceBase{DeviceData: deviceData}
		h.states = map[string]stateItem{
			id: {
				ref:   deviceData.Ref,
				value: state,
			},
		}
		device, err := h.Get(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(device.DeviceData, deviceData) {
			t.Error("expected\n", deviceData, "got\n", device.DeviceData)
		}
		if device.State != state {
			t.Error("expected\n", state, "got\n", device.State)
		}
	})
}

func TestHandler_GetAll(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.DeviceBase)}
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
		stgHdl.devices[id] = lib_model.DeviceBase{DeviceData: deviceData}
		h.states = map[string]stateItem{
			id: {
				ref:   deviceData.Ref,
				value: state,
			},
		}
		devices, err := h.GetAll(context.Background(), lib_model.DevicesFilter{})
		if err != nil {
			t.Error(err)
		}
		if len(devices) != 1 {
			t.Error("expected 1 entry")
		}
		if devices[id].State != state {
			t.Error("expected\n", state, "got\n", devices[id].State)
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

func TestHandler_UpdateUserData(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.DeviceBase)}
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
		if err := h.SetUserData(context.Background(), id, userDataBase); err == nil {
			t.Error("expected error")
		}
		if len(stgHdl.devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.DeviceBase{DeviceData: deviceData}
		if err := h.SetUserData(context.Background(), id, userDataBase); err != nil {
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
		err := h.SetUserData(context.Background(), id, lib_model.DeviceUserDataBase{Attributes: []lib_model.DeviceAttribute{{Value: "test"}}})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_SetStates(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	h := New(nil, 0)
	h.states = map[string]stateItem{
		id: {
			ref: "test",
		},
	}
	if err := h.SetStates(context.Background(), "test", lib_model.Online); err != nil {
		t.Error(err)
	}
	sItem, ok := h.states[id]
	if !ok {
		t.Error("not in map")
	}
	if sItem.ref != "test" {
		t.Error("expected\n", "test", "got\n", sItem.ref)
	}
	if sItem.value != lib_model.Online {
		t.Error("expected\n", lib_model.Online, "got\n", sItem.value)
	}
	t.Run("invalid input", func(t *testing.T) {
		err := h.SetStates(context.Background(), "", "test")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_Delete(t *testing.T) {
	util.InitLogger(util.LoggerConfig{Terminal: true, Level: 4})
	stgHdl := &stgHdlMock{devices: make(map[string]lib_model.DeviceBase)}
	h := New(stgHdl, 0)
	t.Run("does not exist", func(t *testing.T) {
		if err := h.Delete(context.Background(), id); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("exists", func(t *testing.T) {
		stgHdl.devices[id] = lib_model.DeviceBase{DeviceData: deviceData}
		if err := h.Delete(context.Background(), id); err != nil {
			t.Error(err)
		}
		if len(stgHdl.devices) != 0 {
			t.Error("expected 0 entries")
		}
	})
}

func Test_validateDeviceBase(t *testing.T) {
	dData := lib_model.DeviceDataBase{
		ID:   "test",
		Ref:  "test",
		Type: "test",
		Attributes: []lib_model.DeviceAttribute{
			{
				Key:   "test",
				Value: "test",
			},
		},
	}
	t.Run("valid", func(t *testing.T) {
		if err := validateDeviceData(dData); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid ID", func(t *testing.T) {
		idb := dData
		idb.ID = ""
		if err := validateDeviceData(idb); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid type", func(t *testing.T) {
		idd := dData
		idd.Type = ""
		if err := validateDeviceData(idd); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid ref", func(t *testing.T) {
		idd := dData
		idd.Ref = ""
		if err := validateDeviceData(idd); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("invalid attribute", func(t *testing.T) {
		idd := dData
		idd.Attributes = append(idd.Attributes, lib_model.DeviceAttribute{})
		if err := validateDeviceData(idd); err == nil {
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

func Test_validateState(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run(lib_model.Online, func(t *testing.T) {
			if err := validateState(lib_model.Online); err != nil {
				t.Error(err)
			}
		})
		t.Run(lib_model.Offline, func(t *testing.T) {
			if err := validateState(lib_model.Offline); err != nil {
				t.Error(err)
			}
		})
		t.Run("empty", func(t *testing.T) {
			if err := validateState(""); err != nil {
				t.Error(err)
			}
		})
	})
	t.Run("invalid", func(t *testing.T) {
		if err := validateState("test"); err == nil {
			t.Error("expected error")
		}
	})
}

type stgHdlMock struct {
	devices   map[string]lib_model.DeviceBase
	getAllErr error
}

func (m *stgHdlMock) BeginTransaction(_ context.Context) (driver.Tx, error) {
	panic("not implemented")
}

func (m *stgHdlMock) Create(_ context.Context, tx driver.Tx, dBase lib_model.DeviceData) error {
	if tx != nil {
		panic("not implemented")
	}
	if _, ok := m.devices[dBase.ID]; ok {
		return errors.New("duplicate device")
	}
	m.devices[dBase.ID] = lib_model.DeviceBase{DeviceData: dBase}
	return nil
}

func (m *stgHdlMock) Read(_ context.Context, id string) (lib_model.DeviceBase, error) {
	device, ok := m.devices[id]
	if !ok {
		return device, lib_model.NewNotFoundError(errors.New("not found"))
	}
	return device, nil
}

func (m *stgHdlMock) ReadAll(_ context.Context, _ lib_model.DevicesFilter) (map[string]lib_model.DeviceBase, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.devices, nil
}

func (m *stgHdlMock) Update(_ context.Context, tx driver.Tx, dBase lib_model.DeviceData) error {
	if tx != nil {
		panic("not implemented")
	}
	device, ok := m.devices[dBase.ID]
	if !ok {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	device.DeviceData = dBase
	m.devices[dBase.ID] = device
	return nil
}

func (m *stgHdlMock) UpdateUserData(_ context.Context, tx driver.Tx, id string, userData lib_model.DeviceUserData) error {
	if tx != nil {
		panic("not implemented")
	}
	device, ok := m.devices[id]
	if !ok {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	device.UserData = userData
	m.devices[id] = device
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
