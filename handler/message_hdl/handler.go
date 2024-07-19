package message_hdl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util/topic"
)

type Handler struct {
	devicesHdl handler.DevicesHandler
}

func New(devicesHdl handler.DevicesHandler) *Handler {
	return &Handler{devicesHdl: devicesHdl}
}

func (h *Handler) HandleMessage(m handler.Message) error {
	var ref string
	switch {
	case parseTopic(topic.DevicesSub, m.Topic(), &ref):
		var dm lib_model.DeviceMessage
		if err := json.Unmarshal(m.Payload(), &dm); err != nil {
			return err
		}
		switch dm.Method {
		case lib_model.Set:
			if dm.Data == nil {
				return errors.New("missing device data")
			}
			err := h.devicesHdl.Put(context.Background(), lib_model.DeviceData{
				ID:         dm.DeviceID,
				Ref:        ref,
				Name:       dm.Data.Name,
				State:      dm.Data.State,
				Type:       dm.Data.Type,
				Attributes: dm.Data.Attributes,
			})
			if err != nil {
				var iie *lib_model.InvalidInputError
				if errors.As(err, &iie) {
					return err
				}
			}
		case lib_model.Delete:
			_ = h.devicesHdl.Delete(context.Background(), dm.DeviceID)
		default:
			return fmt.Errorf("unknown method '%s'", dm.Method)
		}
	case parseTopic(topic.LastWillSub, m.Topic(), &ref):
		_ = h.devicesHdl.SetStates(context.Background(), ref, lib_model.Offline)
	default:
		return fmt.Errorf("parsing topic '%s' failed", m.Topic())
	}
	return nil
}
