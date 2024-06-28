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

type Handler struct {
	stgHdl  handler.DevicesStorageHandler
	timeout time.Duration
	mu      sync.RWMutex
}

func New(stgHdl handler.DevicesStorageHandler, timeout time.Duration) *Handler {
	return &Handler{
		stgHdl:  stgHdl,
		timeout: timeout,
	}
}

func (h *Handler) Set(ctx context.Context, deviceBase lib_model.DeviceBase) error {
	err := validateDeviceBase(deviceBase)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceBase.ID)
	if err != nil {
		var nfe *lib_model.NotFoundError
		if !errors.As(err, &nfe) {
			return err
		}
		ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
		defer cf2()
		return h.stgHdl.Create(ctxWt2, nil, lib_model.Device{
			DeviceBase: deviceBase,
			Created:    time.Now().UTC(),
		})
	}
	device.DeviceBase = deviceBase
	device.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	return h.stgHdl.Update(ctxWt2, nil, device)
}

func (h *Handler) Add(ctx context.Context, deviceBase lib_model.DeviceBase) error {
	err := validateDeviceBase(deviceBase)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	return h.stgHdl.Create(ctxWt, nil, lib_model.Device{
		DeviceBase: deviceBase,
		Created:    time.Now().UTC(),
	})
}

func (h *Handler) Get(ctx context.Context, id string) (lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	return h.stgHdl.Read(ctxWt, id)
}

func (h *Handler) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	return h.stgHdl.ReadAll(ctxWt, filter)
}

func (h *Handler) Update(ctx context.Context, deviceBase lib_model.DeviceBase) error {
	err := validateDeviceBase(deviceBase)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceBase.ID)
	if err != nil {
		return err
	}
	device.DeviceBase = deviceBase
	device.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	return h.stgHdl.Update(ctxWt2, nil, device)
}

func (h *Handler) UpdateUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	if err := validateAttributes(userDataBase.Attributes); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	if !h.mu.TryLock() {
		return lib_model.NewResourceBusyError(errors.New("acquiring lock failed"))
	}
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		return err
	}
	device.UserData.DeviceUserDataBase = userDataBase
	device.UserData.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	return h.stgHdl.Update(ctxWt2, nil, device)
}

func (h *Handler) Delete(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	return h.stgHdl.Delete(ctxWt, nil, id)
}

func validateDeviceBase(dBase lib_model.DeviceBase) error {
	if dBase.ID == "" {
		return errors.New("empty id")
	}
	if dBase.Type == "" {
		return errors.New("empty type")
	}
	if !isValidDeviceState(dBase.State) {
		return fmt.Errorf("invalid state: %s", dBase.State)
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

func isValidDeviceState(s string) bool {
	switch s {
	case "":
		return true
	case lib_model.Online:
		return true
	case lib_model.Offline:
		return true
	default:
		return false
	}
}
