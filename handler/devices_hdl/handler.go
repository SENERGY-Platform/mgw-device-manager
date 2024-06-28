package devices_hdl

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"sync"
	"time"
)

const logPrefix = "[device-hdl]"

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
	util.Logger.Debugf("set device (%+v)", deviceBase)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceBase.ID)
	if err != nil {
		var nfe *lib_model.NotFoundError
		if !errors.As(err, &nfe) {
			util.Logger.Errorf("set device (%+v): %s", deviceBase, err)
			return err
		}
		ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
		defer cf2()
		err = h.stgHdl.Create(ctxWt2, nil, lib_model.Device{
			DeviceBase: deviceBase,
			Created:    time.Now().UTC(),
		})
		if err != nil {
			util.Logger.Errorf("set device (%+v): %s", deviceBase, err)
			return err
		}
		return nil
	}
	device.DeviceBase = deviceBase
	device.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	if err = h.stgHdl.Update(ctxWt2, nil, device); err != nil {
		util.Logger.Errorf("set device (%+v): %s", deviceBase, err)
		return err
	}
	return nil
}

func (h *Handler) Add(ctx context.Context, deviceBase lib_model.DeviceBase) error {
	err := validateDeviceBase(deviceBase)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("add device (%+v)", deviceBase)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	err = h.stgHdl.Create(ctxWt, nil, lib_model.Device{
		DeviceBase: deviceBase,
		Created:    time.Now().UTC(),
	})
	if err != nil {
		util.Logger.Errorf("add device (%+v): %s", deviceBase, err)
		return err
	}
	return nil
}

func (h *Handler) Get(ctx context.Context, id string) (lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	util.Logger.Debugf("get device (%s)", id)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		util.Logger.Errorf("get device (%s): %s", id, err)
		return lib_model.Device{}, err
	}
	return device, nil
}

func (h *Handler) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	util.Logger.Debugf("get devices (%+v)", filter)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	devices, err := h.stgHdl.ReadAll(ctxWt, filter)
	if err != nil {
		util.Logger.Errorf("get devices (%+v): %s", filter, err)
		return nil, err
	}
	return devices, nil
}

func (h *Handler) Update(ctx context.Context, deviceBase lib_model.DeviceBase) error {
	err := validateDeviceBase(deviceBase)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("update device (%+v)", deviceBase)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceBase.ID)
	if err != nil {
		util.Logger.Errorf("update device (%+v): %s", deviceBase, err)
		return err
	}
	device.DeviceBase = deviceBase
	device.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	if err = h.stgHdl.Update(ctxWt2, nil, device); err != nil {
		util.Logger.Errorf("update device (%+v): %s", deviceBase, err)
		return err
	}
	return nil
}

func (h *Handler) UpdateUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	if err := validateAttributes(userDataBase.Attributes); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("update device user data (%+v)", userDataBase)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		util.Logger.Errorf("update device user data (%+v): %s", userDataBase, err)
		return err
	}
	device.UserData.DeviceUserDataBase = userDataBase
	device.UserData.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	if err = h.stgHdl.Update(ctxWt2, nil, device); err != nil {
		util.Logger.Errorf("update device user data (%+v): %s", userDataBase, err)
		return err
	}
	return nil
}

func (h *Handler) Delete(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("delete device (%s)", id)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	if err := h.stgHdl.Delete(ctxWt, nil, id); err != nil {
		util.Logger.Errorf("delete device (%s): %s", id, err)
		return err
	}
	return nil
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
