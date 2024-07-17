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

func (h *Handler) Put(ctx context.Context, deviceData lib_model.DeviceData) error {
	err := validateDeviceData(deviceData)
	if err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("%s put device (%+v)", logPrefix, deviceData)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, deviceData.ID)
	if err != nil {
		var nfe *lib_model.NotFoundError
		if !errors.As(err, &nfe) {
			util.Logger.Errorf("%s put device (%+v): %s", logPrefix, deviceData, err)
			return err
		}
		ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
		defer cf2()
		err = h.stgHdl.Create(ctxWt2, nil, lib_model.DeviceBase{
			DeviceData: deviceData,
			Created:    time.Now().UTC(),
		})
		if err != nil {
			util.Logger.Errorf("%s put device (%+v): %s", logPrefix, deviceData, err)
			return err
		}
		return nil
	}
	device.DeviceData = deviceData
	device.Updated = time.Now().UTC()
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	if err = h.stgHdl.Update(ctxWt2, nil, device.DeviceBase); err != nil {
		util.Logger.Errorf("%s put device (%+v): %s", logPrefix, deviceData, err)
		return err
	}
	return nil
}

func (h *Handler) Get(ctx context.Context, id string) (lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	util.Logger.Debugf("%s get device (%s)", logPrefix, id)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	device, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		util.Logger.Errorf("%s get device (%s): %s", logPrefix, id, err)
		return lib_model.Device{}, err
	}
	return device, nil
}

func (h *Handler) GetAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	util.Logger.Debugf("%s get devices (%+v)", logPrefix, filter)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	devices, err := h.stgHdl.ReadAll(ctxWt, filter)
	if err != nil {
		util.Logger.Errorf("%s get devices (%+v): %s", logPrefix, filter, err)
		return nil, err
	}
	return devices, nil
}

func (h *Handler) SetUserData(ctx context.Context, id string, userDataBase lib_model.DeviceUserDataBase) error {
	if err := validateAttributes(userDataBase.Attributes); err != nil {
		return lib_model.NewInvalidInputError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("%s set device user data (%+v)", logPrefix, userDataBase)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	_, err := h.stgHdl.Read(ctxWt, id)
	if err != nil {
		util.Logger.Errorf("%s set device user data (%+v): %s", logPrefix, userDataBase, err)
		return err
	}
	ctxWt2, cf2 := context.WithTimeout(ctx, h.timeout)
	defer cf2()
	err = h.stgHdl.UpdateUserData(ctxWt2, nil, id, lib_model.DeviceUserData{
		DeviceUserDataBase: userDataBase,
		Updated:            time.Now().UTC(),
	})
	if err != nil {
		util.Logger.Errorf("%s set device user data (%+v): %s", logPrefix, userDataBase, err)
		return err
	}
	return nil
}

func (h *Handler) SetStates(ctx context.Context, ref string, state lib_model.DeviceState) error {
	if !isValidDeviceState(state) {
		return lib_model.NewInvalidInputError(errors.New("invalid state"))
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("%s set device sates (%s %s)", logPrefix, ref, state)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	if err := h.stgHdl.UpdateStates(ctxWt, nil, ref, state, time.Now().UTC()); err != nil {
		util.Logger.Errorf("%s set device sates (%s %s): %s", logPrefix, ref, state, err)
		return err
	}
	return nil
}

func (h *Handler) Delete(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	util.Logger.Debugf("%s delete device (%s)", logPrefix, id)
	ctxWt, cf := context.WithTimeout(ctx, h.timeout)
	defer cf()
	if err := h.stgHdl.Delete(ctxWt, nil, id); err != nil {
		util.Logger.Errorf("%s delete device (%s): %s", logPrefix, id, err)
		return err
	}
	return nil
}

func validateDeviceData(dBase lib_model.DeviceData) error {
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
