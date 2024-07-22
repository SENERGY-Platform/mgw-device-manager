package devices_hdl

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"sync"
	"time"
)

type stateItem struct {
	ref   string
	value lib_model.DeviceState
}

type Handler struct {
	stgHdl  handler.DevicesStorageHandler
	timeout time.Duration
	states  map[string]stateItem
	mu      sync.RWMutex
}

func New(stgHdl handler.DevicesStorageHandler, timeout time.Duration) *Handler {
	return &Handler{
		stgHdl:  stgHdl,
		timeout: timeout,
		states:  make(map[string]stateItem),
	}
}

func (h *Handler) Put(ctx context.Context, deviceData lib_model.DeviceDataBase, state lib_model.DeviceState) error {
	if err := validateDeviceData(deviceData); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	if err := validateState(state); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceData.ID)
	if err != nil {
		var nfe *lib_model.NotFoundError
		if !errors.As(err, &nfe) {
			return fmt.Errorf("put device: %s", err)
		}
		ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
		defer cf2()
		err = h.stgHdl.Create(ctxWt2, nil, lib_model.DeviceData{
			DeviceDataBase: deviceData,
			Created:        time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("put device: %s", err)
		}
	} else {
		device.DeviceDataBase = deviceData
		device.Updated = time.Now().UTC()
		ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
		defer cf2()
		if err = h.stgHdl.Update(ctxWt2, nil, device.DeviceData); err != nil {
			return fmt.Errorf("put device: %s", err)
		}
	}
	h.states[deviceData.ID] = stateItem{
		ref:   deviceData.Ref,
		value: state,
	}
	return nil
}

func (h *Handler) Get(ctx context.Context, id string) (lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		return lib_model.Device{}, fmt.Errorf("get device: %s", err)
	}
	return lib_model.Device{
		DeviceBase: device,
		State:      h.getState(id),
	}, nil
}

func (h *Handler) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	deviceBases, err := h.stgHdl.ReadAll(ctxWt, filter)
	if err != nil {
		return nil, fmt.Errorf("get devices: %s", err)
	}
	devices := make(map[string]lib_model.Device)
	for id, deviceBase := range deviceBases {
		state := h.getState(id)
		if filter.State != "" && state != filter.State {
			continue
		}
		devices[id] = lib_model.Device{
			DeviceBase: deviceBase,
			State:      state,
		}
	}
	return devices, nil
}

func (h *Handler) SetUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	if err := validateAttributes(userDataBase.Attributes); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	_, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		return fmt.Errorf("set device user data: %s", err)
	}
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	err = h.stgHdl.UpdateUserData(ctxWt2, nil, id, lib_model.DeviceUserData{
		DeviceUserDataBase: userDataBase,
		Updated:            time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("set device user data: %s", err)
	}
	return nil
}

func (h *Handler) SetStates(ctx context.Context, ref string, state lib_model.DeviceState) error {
	if err := validateState(state); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, sItem := range h.states {
		if sItem.ref == ref {
			sItem.value = state
			h.states[id] = sItem
		}
	}
	return nil
}

func (h *Handler) Delete(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	if err := h.stgHdl.Delete(ctxWt, nil, id); err != nil {
		return fmt.Errorf("delete device: %s", err)
	}
	delete(h.states, id)
	return nil
}

func (h *Handler) getState(id string) string {
	sItem, ok := h.states[id]
	if !ok || sItem.value == "" {
		return lib_model.NotAvailable
	}
	return sItem.value
}

func validateDeviceData(dBase lib_model.DeviceDataBase) error {
	if dBase.ID == "" {
		return errors.New("empty id")
	}
	if dBase.Type == "" {
		return errors.New("empty type")
	}
	if dBase.Ref == "" {
		return errors.New("empty reference")
	}
	return validateAttributes(dBase.Attributes)
}

func validateAttributes(attrs []lib_model.DeviceAttribute) error {
	for _, attr := range attrs {
		if attr.Key == "" {
			return errors.New("empty attribute key")
		}
	}
	return nil
}

func validateState(s string) error {
	switch s {
	case "":
		return nil
	case lib_model.Online:
		return nil
	case lib_model.Offline:
		return nil
	default:
		return errors.New("invalid state")
	}
}
