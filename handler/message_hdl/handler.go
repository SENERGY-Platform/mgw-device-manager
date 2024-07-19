package message_hdl

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"github.com/SENERGY-Platform/mgw-device-manager/util/topic"
)

const logPrefix = "[message-hdl]"

type Handler struct {
	devicesHdl handler.DevicesHandler
}

func New(devicesHdl handler.DevicesHandler) *Handler {
	return &Handler{devicesHdl: devicesHdl}
}

func (h *Handler) HandleMessage(m handler.Message) {
	var ref string
	switch {
	case parseTopic(topic.DevicesSub, m.Topic(), &ref):
		var dm lib_model.DeviceMessage
		if err := json.Unmarshal(m.Payload(), &dm); err != nil {
			util.Logger.Errorf("%s unmarshal message: %s", logPrefix, err)
			return
		}
		switch dm.Method {
		case lib_model.Set:
			util.Logger.Infof("%s set device (%s)", logPrefix, dm.DeviceID)
			if dm.Data == nil {
				util.Logger.Errorf("%s set device (%s): missing data", logPrefix, dm.DeviceID)
				return
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
				util.Logger.Errorf("%s set device (%s): %s", logPrefix, dm.DeviceID, err)
			}
		case lib_model.Delete:
			util.Logger.Infof("%s delete device (%s)", logPrefix, dm.DeviceID)
			if err := h.devicesHdl.Delete(context.Background(), dm.DeviceID); err != nil {
				util.Logger.Errorf("%s delete device (%s): %s", logPrefix, dm.DeviceID, err)
			}
		default:
			util.Logger.Errorf("%s unknown method '%s'", logPrefix, dm.Method)
		}
	case parseTopic(topic.LastWillSub, m.Topic(), &ref):
		util.Logger.Infof("%s set device states (%s)", logPrefix, ref)
		if err := h.devicesHdl.SetStates(context.Background(), ref, lib_model.Offline); err != nil {
			util.Logger.Errorf("%s set device states (%s): %s", logPrefix, ref, err)
		}
	default:
		util.Logger.Errorf("%s unknown topic '%s'", logPrefix, m.Topic())
	}
	return
}
