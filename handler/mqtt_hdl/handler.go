package mqtt_hdl

import (
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"github.com/SENERGY-Platform/mgw-device-manager/util/topic"
)

const LogPrefix = "[mqtt-hdl]"

const (
	RelayMsgErrString  = "%s relay message (%s): %s"
	SubscribeString    = "%s subscribe topic (%s)"
	SubscribedString   = "%s subscribed topic (%s)"
	SubscribeErrString = "%s subscribe topic (%s): %s"
)

type Handler struct {
	client          handler.MqttClient
	qos             byte
	messageRelayHdl handler.MessageRelayHandler
}

func New(qos byte, messageRelayHdl handler.MessageRelayHandler) *Handler {
	return &Handler{
		qos:             qos,
		messageRelayHdl: messageRelayHdl,
	}
}

func (h *Handler) SetMqttClient(c handler.MqttClient) {
	h.client = c
}

func (h *Handler) HandleOnConnect() {
	if err := h.handleSubscriptions(); err == nil {
		go h.publishRefreshSignal()
	}
}

func (h *Handler) handleSubscriptions() error {
	util.Logger.Debugf(SubscribeString, LogPrefix, topic.DevicesSub)
	err := h.client.Subscribe(topic.DevicesSub, h.qos, func(m handler.Message) {
		if err := h.messageRelayHdl.Put(m); err != nil {
			util.Logger.Errorf(RelayMsgErrString, LogPrefix, m.Topic(), err)
		}
	})
	if err != nil {
		util.Logger.Errorf(SubscribeErrString, LogPrefix, topic.DevicesSub, err)
		return err
	}
	util.Logger.Infof(SubscribedString, LogPrefix, topic.DevicesSub)
	util.Logger.Debugf(SubscribeString, LogPrefix, topic.LastWillSub)
	err = h.client.Subscribe(topic.LastWillSub, h.qos, func(m handler.Message) {
		if err := h.messageRelayHdl.Put(m); err != nil {
			util.Logger.Errorf(RelayMsgErrString, LogPrefix, m.Topic(), err)
		}
	})
	if err != nil {
		util.Logger.Errorf(SubscribeErrString, LogPrefix, topic.LastWillSub, err)
		return err
	}
	util.Logger.Infof(SubscribedString, LogPrefix, topic.LastWillSub)
	return nil
}

func (h *Handler) publishRefreshSignal() {
	if err := h.client.Publish(topic.RefreshPub, h.qos, false, []byte("1")); err != nil {
		util.Logger.Errorf("publish refresh signal: %s", err)
	}
}
